# Sécurité — ProfilsActifs

Ce document décrit les mesures de sécurité réellement implémentées, l'endroit
où elles se trouvent dans le code, et les limitations connues. Rien n'est
théorique : chaque point est vérifiable dans les fichiers cités.

## Périmètre

- API HTTP Go (`backend/src/`), base PostgreSQL via GORM, front Vue servi séparément.
- Authentification par **JWT porté dans l'en-tête `Authorization`** (pas de cookie de session).
- Le catalogue de profils et les fiches sont **publics** (exigence fonctionnelle).
- Le chiffrement du transport (TLS) est délégué à un reverse-proxy en amont.

---

## 1. Mots de passe

| Quoi | Où |
| --- | --- |
| Hachage **bcrypt** avec sel automatique par hash, coût par défaut (10) | `backend/src/token.go` — `hashPassword` / `checkPassword` |
| Le hash n'est **jamais** sérialisé dans les réponses JSON | `backend/models/user.go` — `PasswordHash string \`json:"-"\`` |

- Le mot de passe en clair n'est jamais stocké ni journalisé : il est haché dès la réception dans `register`.
- `checkPassword` s'appuie sur `bcrypt.CompareHashAndPassword`, résistant au timing par construction.
- `PasswordHash` étant marqué `json:"-"`, il est absent de `GET /me` et de `GET /users` même si l'objet `User` complet est renvoyé.

**À dire :** « Les mots de passe sont hachés en bcrypt, avec un sel généré
automatiquement et unique par utilisateur. Le champ n'est jamais sérialisé,
donc jamais renvoyé au client. »

## 2. Jetons JWT

| Quoi | Où |
| --- | --- |
| Signature **HS256**, secret hors du code (`JWT_SECRET`) | `backend/src/token.go` — `generateToken` |
| **Rejet explicite** de tout jeton dont l'algorithme n'est pas HMAC | `backend/src/auth.go` — `parseToken` (`t.Method.(*jwt.SigningMethodHMAC)`) |
| **Expiration à 24 h** (`exp`), vérifiée automatiquement par `jwt/v5` | `backend/src/token.go` |

- Le contrôle `SigningMethodHMAC` est la parade classique à l'attaque `alg=none`
  et à la confusion de clé (un attaquant qui passerait un jeton signé RS256 avec
  la clé publique est rejeté).
- Aucune session serveur : le jeton est auto-porté, à durée de vie limitée.
- Le secret provient de `os.Getenv("JWT_SECRET")`, fourni par `.env` (non versionné).

**À dire :** « Le jeton est signé en HMAC-SHA256. À la vérification, on impose
que la méthode de signature soit HMAC : un `alg=none` ou une confusion RSA/HMAC
est rejetée. Durée de vie 24 h, pas de session côté serveur. »

## 3. Contrôle d'accès (RBAC)

| Quoi | Où |
| --- | --- |
| Middleware `requireAuth` — jeton valide obligatoire | `backend/src/auth.go` |
| Middleware `requireRole("candidate", …)` — rôle exact exigé | `backend/src/auth.go` |
| Application route par route | `backend/src/main.go` |

Routes protégées :

- `requireRole("candidate")` : `POST /videos`, `POST /questionnaire/start|answer|validate`
- `requireAuth` : `GET /me`, `GET|POST|DELETE /consent`
- L'identifiant utilisateur est injecté dans le contexte de la requête depuis le
  jeton (`currentUserID`) — **jamais lu depuis le corps de la requête**. Un
  utilisateur ne peut donc agir que sur ses propres données (vidéos, réponses,
  consentement).

**Pas d'élévation de privilège à l'inscription :**

- `register` décode un struct anonyme de **4 champs seulement**
  (`name`, `email`, `password`, `birthday`) — pas de *mass assignment*.
- Le rôle est forcé en dur à `"candidate"` (`backend/src/auth.go`, création du `models.User`).
- Les rôles `recruiter` / `admin` existent dans le modèle mais aucun endpoint ni seed ne permet de les obtenir.

**À dire :** « Chaque route sensible passe par un middleware qui vérifie le jeton
et, si besoin, le rôle. L'ID de l'utilisateur vient toujours du jeton, jamais du
payload, donc pas d'accès horizontal aux données d'autrui. Le rôle est imposé
côté serveur à l'inscription. »

## 4. Validation des entrées

| Contrôle | Où |
| --- | --- |
| Champs obligatoires, mot de passe ≥ 6 caractères | `backend/src/auth.go` — `register` |
| Format de date `YYYY-MM-DD` strict | `backend/src/auth.go` — `time.Parse` |
| **Âge minimum 16 ans** → `403` sinon | `backend/src/token.go` `computeAge` + `backend/src/auth.go` |
| Statut « jeune » (16–17 ans) calculé côté serveur | `backend/src/auth.go` (`models.StatusYouth`) |
| E-mail normalisé en minuscules + **unicité** en base | `backend/src/auth.go`, `backend/models/user.go` (`uniqueIndex`) |
| Questionnaire : `question_id` existe, `choice` ∈ options de la question | `backend/src/quiz.go` — `submitAnswer` |
| Une seule réponse par (utilisateur, question) | `backend/models/quiz.go` — `uniqueIndex:idx_user_question` |

**À dire :** « On refuse l'inscription des moins de 16 ans, on contrôle le format
de la date et les champs obligatoires. Sur le questionnaire, on vérifie que la
question existe et que la réponse fait partie des options prévues, et un index
unique empêche de répondre deux fois. »

## 5. Injection SQL

- **100 % des accès base passent par GORM** avec des requêtes paramétrées :
  `DB.Where("email = ?", …)`, `DB.First(&user, id)`, `DB.Model(...).Update(...)`.
- Aucune concaténation de chaîne dans une requête, nulle part dans le code.
- Les identifiants d'URL (`GET /users/{id}`) sont convertis en entier
  (`strconv.Atoi`) avant toute utilisation (`backend/src/users.go`).

**À dire :** « Aucune requête n'est construite par concaténation. L'ORM
paramètre systématiquement les valeurs, y compris les identifiants d'URL qui
sont d'abord convertis en entier. »

## 6. Upload de fichiers

| Quoi | Où |
| --- | --- |
| Taille limitée à **100 Mo** (`http.MaxBytesReader` + `ParseMultipartForm`) | `backend/src/videos.go` — `maxUploadSize` |
| **Whitelist d'extensions** : `.mp4`, `.webm`, `.mov` uniquement | `backend/src/videos.go` — `uploadFile` |
| Nom de fichier **généré côté serveur** (`<horodatage nanoseconde>.<ext>`) | `backend/src/videos.go` |
| Fichiers servis en lecture seule depuis `uploads/` | `backend/src/main.go` — `http.FileServer` |

- Le nom du fichier envoyé par le client n'est **jamais** réutilisé pour écrire
  sur le disque : on ne garde que l'extension, le nom final est généré. Pas de
  traversée de répertoire possible.
- `http.MaxBytesReader` coupe la lecture du corps au-delà de la limite (pas de
  bufférisation illimitée en mémoire).

**À dire :** « L'upload est plafonné à 100 Mo via `MaxBytesReader`, seules trois
extensions vidéo sont acceptées, et le nom de fichier est régénéré côté serveur
pour éviter toute traversée de chemin. »

## 7. RGPD — consentement et droit à l'effacement

| Quoi | Où |
| --- | --- |
| Consentement **explicite, versionné et horodaté** | `backend/src/consent.go`, `backend/models/consent.go` (`Version`, `ApprouveAt`) |
| Révocation = **soft-delete** (`revoked_at`) → trace conservée | `backend/src/consent.go` — `revokeConsent` |
| La révocation **supprime les vidéos** : lignes en base **et** fichiers sur disque | `backend/src/consent.go` — `deleteUserVideos` (`os.Remove`) |
| Suppression d'un utilisateur → cascade sur vidéos, réponses, consentements, badge | `backend/models/*.go` — `constraint:OnDelete:CASCADE` |

**À dire :** « Le consentement est tracé avec sa version et sa date. La
révocation ne fait pas que masquer : elle marque le consentement comme révoqué
(on garde la preuve) et déclenche la suppression physique des vidéos, fichiers
compris. Les relations sont en `ON DELETE CASCADE`. »

## 8. Gestion des secrets

| Quoi | Où |
| --- | --- |
| Aucun secret versionné : seul `.env.example` (valeurs vides) est suivi | `.env.example` |
| `.env` est ignoré par Git | `.gitignore` |
| Les secrets sont injectés à l'exécution par variables d'environnement | `docker-compose.yml` (`env_file`, `${JWT_SECRET}`) |

**À dire :** « Aucun secret dans le dépôt. On versionne un `.env.example` avec
des champs vides, le vrai `.env` est gitignoré et les valeurs sont passées par
l'environnement au démarrage. »

## 9. Conteneurs

| Quoi | Où |
| --- | --- |
| Build **multi-stage** : image de build ≠ image d'exécution | `backend/Dockerfile` |
| Binaire **statique** (`CGO_ENABLED=0`) sur `alpine:3.20` minimal | `backend/Dockerfile` |
| `.dockerignore` limite le contexte de build | `backend/.dockerignore` |

**À dire :** « Build en deux étapes : on compile dans une image Go complète, on
n'embarque que le binaire statique dans une Alpine minimale. Surface d'attaque
réduite, pas de toolchain ni de sources dans l'image finale. »

## 10. Transport et en-têtes

- Authentification par en-tête `Authorization: Bearer` — **pas de cookie**, donc
  **pas de vecteur CSRF** (rien n'est envoyé automatiquement par le navigateur).
- CORS ouvert (`Access-Control-Allow-Origin: *`, `backend/src/auth.go` —
  `corsMiddleware`) : acceptable dans ce modèle (jeton en en-tête), à restreindre
  à l'origine du front en production.
- Vue échappe le HTML par défaut ; aucune donnée utilisateur n'est rendue via
  `v-html`.
