# Tests — ProfilsActifs

## 1. Tests automatisés (backend)

Fichier unique : **`backend/src/main_test.go`** — 5 tests unitaires, sans base de
données (uniquement bcrypt, JWT et calcul d'âge).

| Test | Fonction visée | Ce qu'il vérifie |
| --- | --- | --- |
| `TestPassword` | `hashPassword`, `checkPassword` (`token.go`) | le hash bcrypt est différent du mot de passe en clair ; le bon mot de passe est accepté ; un mauvais est refusé |
| `TestComputeAge` | `computeAge` (`token.go`) | une date de naissance à −20 ans donne bien `20` |
| `TestToken` | `generateToken` + `parseToken` (`token.go`, `auth.go`) | un jeton généré puis relu redonne le bon `id` (7) et le bon rôle (`candidate`) |
| `TestTokenManquant` | `parseToken` | une requête sans en-tête `Authorization` renvoie une erreur |
| `TestTokenFalsifie` | `parseToken` | un jeton signé avec un secret, vérifié avec un autre secret, est **rejeté** |

Ces tests couvrent la partie sensible : hachage des mots de passe et
vérification des jetons (voir [SECURITE.md](SECURITE.md) §1 et §2).

### Couverture

```
go test ./src/... -cover
```

| Fonction | Couverture |
| --- | --- |
| `hashPassword` | 100 % |
| `checkPassword` | 100 % |
| `generateToken` | 100 % |
| `parseToken` | 85 % |
| `computeAge` | 80 % |
| **global (`profilsactifs/src`)** | **7,9 %** |

La couverture globale est faible parce que les *handlers* HTTP (`register`,
`login`, `uploadVideo`, `consent`, `quiz`) et la couche base ne sont pas testés
unitairement — ils sont vérifiés manuellement (§3).

## 2. Lancer les tests

### En local

```bash
cd backend
go test ./... -v          # verbeux
go test ./... -cover      # avec le taux de couverture
go vet ./...              # analyse statique
go build ./...            # compilation
```

Aucune base n'est nécessaire pour `go test` (les 5 tests sont hors DB).

### En intégration continue

`.github/workflows/ci.yml`, job `backend`, à chaque `push` et `pull_request` :

```
go vet ./...
go build ./...
go test ./... -v
```

> Le fichier de test **doit** s'appeler `*_test.go` (`main_test.go`) pour que
> `go test` le prenne en compte. Nommé `test.go`, il serait compilé mais jamais
> exécuté, et la CI afficherait `[no test files]` en restant verte à tort.

Le job `frontend` du même workflow lance `npm ci`, `vue-tsc -b` (typage) et
`npm run build` — il n'y a pas de test unitaire front, seulement la
vérification que l'application compile.

## 3. Procédure de test manuel (bout en bout)

Prérequis : `cp .env.example .env` (remplir les valeurs), puis
`docker compose up --build`. API sur `http://localhost:8080`, interface sur
`http://localhost:5173`.

### 3.1 Santé de l'API

```bash
curl -s http://localhost:8080/health        # -> {"status":"ok"}
```

### 3.2 Inscription

| Cas | Requête | Attendu |
| --- | --- | --- |
| Nominal (adulte) | `POST /register` avec `birthday` = date ≥ 18 ans | `201`, `{token, user{role:"candidate", status:"adult"}}` |
| Jeune (16–17 ans) | `birthday` entre 16 et 18 ans | `201`, `status:"youth"` |
| Mineur < 16 ans | `birthday` < 16 ans | `403` « you must be at least 16 years old » |
| Mot de passe court | `password` < 6 caractères | `400` |
| Date invalide | `birthday` = `"12/03/2000"` | `400` « invalid birthday » |
| E-mail déjà pris | rejouer une inscription avec le même e-mail | `409` « email already in use » |

```bash
curl -s -X POST http://localhost:8080/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Léa","email":"lea@test.fr","password":"secret1","birthday":"2000-05-01"}'
```

### 3.3 Connexion

| Cas | Attendu |
| --- | --- |
| Bons identifiants | `200`, `{token, user}` |
| Mauvais mot de passe **ou** e-mail inconnu | `401` « invalid credentials » (message identique dans les deux cas → pas d'énumération) |

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"lea@test.fr","password":"secret1"}' | jq -r .token)
```

### 3.4 Questionnaire → badge

```bash
curl -s -X POST http://localhost:8080/questionnaire/start  -H "Authorization: Bearer $TOKEN"
curl -s http://localhost:8080/questions | jq '.[].id'      # récupérer les IDs

# répondre à chaque question (ici "Oui" partout)
curl -s -X POST http://localhost:8080/questionnaire/answer \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"question_id":1,"choice":"Oui"}'

curl -s -X POST http://localhost:8080/questionnaire/validate -H "Authorization: Bearer $TOKEN"
# -> {"score":..,"badge":true|false}   badge = true si score >= 60 % du total
```

| Cas | Attendu |
| --- | --- |
| `choice` hors des options de la question | `400` « invalid choice » |
| `question_id` inexistant | `400` « unknown question » |
| Répondre 2× à la même question | `400` (index unique `user_id + question_id`) |
| Route appelée sans jeton | `401` |
| Route appelée avec un rôle ≠ `candidate` | `403` |

### 3.5 Vidéo + consentement

```bash
# consentement
curl -s -X POST   http://localhost:8080/consent -H "Authorization: Bearer $TOKEN"   # 201
curl -s           http://localhost:8080/consent -H "Authorization: Bearer $TOKEN"   # {"active":true,...}

# vidéo par lien
curl -s -X POST http://localhost:8080/videos -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' -d '{"url":"https://youtu.be/xxxx"}'

# vidéo par fichier
curl -s -X POST http://localhost:8080/videos -H "Authorization: Bearer $TOKEN" \
  -F file=@presentation.mp4
```

| Cas | Attendu |
| --- | --- |
| Fichier > 100 Mo | `413` « file too big » |
| Extension ≠ mp4/webm/mov | `400` « unsupported format » |
| Upload sans jeton candidat | `401` / `403` |

### 3.6 Révocation du consentement (droit à l'effacement)

```bash
curl -s -X DELETE http://localhost:8080/consent -H "Authorization: Bearer $TOKEN"   # 204
```

Vérifications attendues :

- la ou les vidéos de l'utilisateur ont disparu de la base ;
- les fichiers correspondants ont été **supprimés du disque** (`backend/uploads/`) ;
- `GET /consent` renvoie `{"active":false}` ;
- l'enregistrement de consentement reste en base avec `revoked_at` renseigné
  (trace conservée).

### 3.7 Parcours interface

1. `/signup` — créer un compte (tester une date < 16 ans → message d'erreur affiché).
2. Redirection vers `/quiz` — répondre aux 20 questions → écran de résultat avec ou sans badge.
3. `/ma-video` — cocher le consentement, coller un lien YouTube → la vidéo s'affiche ; téléverser un fichier → lecture dans le lecteur.
4. `/candidats` — la liste est visible **sans être connecté** ; `/candidats/:id` affiche la fiche et la vidéo.
5. `/mon-espace` — révoquer le consentement → la vidéo n'est plus visible.
6. Se déconnecter → les routes `/quiz`, `/ma-video`, `/mon-espace` redirigent vers `/login`.
