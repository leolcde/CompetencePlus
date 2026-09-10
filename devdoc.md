# ProfilsActifs — documentation développeur

## Stack

| Couche | Technologies |
| --- | --- |
| Backend | Go 1.26 · `net/http` (routeur standard) · GORM · JWT (`golang-jwt/v5`) · bcrypt |
| Frontend | Vue 3 · Vite · TypeScript · Tailwind CSS · Vue Router |
| Base | PostgreSQL 16 |
| Infra | Docker Compose |

## Structure

```
.
├── backend/
│   ├── src/            handlers + routeur (package main)
│   │   ├── main.go     déclaration des routes, middleware CORS, listen :8080
│   │   ├── auth.go     register, login, parseToken, requireAuth, requireRole
│   │   ├── users.go    Me, listUsers, getUser
│   │   ├── quiz.go     listQuestions, start / answer / validate
│   │   ├── videos.go   upload (lien JSON ou fichier multipart)
│   │   ├── consent.go  get / grant / revoke (révoquer supprime les vidéos)
│   │   ├── token.go    hash bcrypt, génération JWT, calcul d'âge
│   │   └── store.go    variable globale DB
│   ├── models/         structs GORM : User, Question, Answer, BadgeResult, Video, Consent
│   └── db/             Connect, AutoMigrate, SeedQuestions
├── frontend/
│   └── src/
│       ├── lib/
│       │   ├── api.ts  client HTTP typé (préfixe /api)
│       │   └── auth.ts état de session (token + user en localStorage)
│       ├── router/     routes + garde meta.auth
│       ├── type.ts     types partagés avec l'API
│       ├── components/ AppLayout, LegalBanner, VideoPlayer
│       └── views/      pages
├── docker-compose.yml
└── .env.example
```

## Variables d'environnement

| Variable | Rôle |
| --- | --- |
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` | identifiants PostgreSQL |
| `JWT_SECRET` | secret de signature des tokens JWT (HS256) |
| `POSTGRES_URL` | DSN complet pour un lancement **local** du backend (`host = localhost`) |
| `DB_URL` | DSN utilisé **en Docker** — construit par `docker-compose` à partir des `POSTGRES_*` (`host = postgres`) |

Le backend lit `DB_URL` en priorité, sinon `POSTGRES_URL`.

## Base de données

Pas de migrations SQL. Au démarrage, le backend :

1. se connecte (GORM),
2. exécute `AutoMigrate` sur les 5 modèles,
3. `SeedQuestions` insère les 20 questions si la table est vide.

Repartir de zéro : `docker compose down -v`.

## Lancement manuel

**Base** — un PostgreSQL local correspondant à `POSTGRES_URL`.

```bash
# Backend
cd backend
go run ./src            # écoute sur :8080

# Frontend
cd frontend
npm install
npm run dev             # http://localhost:5173
```

## Build de production

```bash
cd backend  && go build -o server ./src
cd frontend && npm run build          # sortie dans frontend/dist
```

## API

Réponses en JSON (les erreurs sont renvoyées en texte brut).
Les routes protégées attendent `Authorization: Bearer <token>`.

| Méthode | Route | Auth | Description |
| --- | --- | --- | --- |
| `GET` | `/health` | — | état de l'API + ping DB |
| `POST` | `/register` | — | `{ name, email, password, birthday }` → `{ token, user }` — refus si < 16 ans |
| `POST` | `/login` | — | `{ email, password }` → `{ token, user }` |
| `GET` | `/me` | token | profil de l'utilisateur connecté |
| `GET` | `/users` | — | liste des utilisateurs |
| `GET` | `/users/{id}` | — | un utilisateur |
| `GET` | `/questions` | — | les 20 questions du questionnaire |
| `POST` | `/questionnaire/start` | candidate | réinitialise les réponses |
| `POST` | `/questionnaire/answer` | candidate | `{ question_id, choice }` |
| `POST` | `/questionnaire/validate` | candidate | calcule le score → `{ user_id, score, badge }` |
| `POST` | `/videos` | candidate | JSON `{ url }` **ou** `multipart/form-data` (champ `file`) |
| `GET` | `/uploads/<fichier>` | — | vidéos téléversées |
| `GET` / `POST` / `DELETE` | `/consent` | token | consentement RGPD — `DELETE` révoque **et supprime les vidéos** |

### Authentification

- JWT **HS256**, expiration **24 h**, claims `{ sub: userID, role }`.
- `parseToken` rejette tout token dont l'algorithme n'est pas HMAC.
- Rôles : `candidate`, `recruiter`, `admin`. L'inscription force `candidate`.
- Statut `youth` (< 18 ans) / `adult`, dérivé de la date de naissance.
- Mots de passe : bcrypt, jamais renvoyés (`json:"-"`).

### Upload de vidéos

- Taille max **100 Mo** (`http.MaxBytesReader`).
- Extensions autorisées : `.mp4`, `.webm`, `.mov`.
- Fichiers stockés dans `backend/uploads/`, servis sur `/uploads/`.

## Frontend

SPA Vue 3. La session (token + user) est persistée en `localStorage`.
Tous les appels réseau passent par `src/lib/api.ts`.

| Route | Accès | Vue |
| --- | --- | --- |
| `/` | public | `HomeView` |
| `/login` · `/signup` | public | `LoginView` · `SignupView` |
| `/candidats` · `/candidats/:id` | public | `CandidatesView` · `CandidateView` |
| `/quiz` | connecté | `QuizView` |
| `/ma-video` | connecté | `MyVideoView` |
| `/mon-espace` | connecté | `AccountView` |
| `*` | public | `NotFoundView` |

La garde `router.beforeEach` redirige vers `/login` toute route `meta.auth` sans session.

### Préfixe `/api`

Le client appelle `/api/login`, `/api/users`, etc. Le proxy Vite
(`vite.config.ts`) réécrit `/api/*` → backend et **retire** le préfixe.
Cela sépare l'espace des **pages** (`/login`) de l'espace des **appels API**
(`/api/login`) : une route de page ne peut jamais entrer en conflit avec une route
du backend, quelle que soit celle qu'on ajoute.

### Bandeau légal

Le composant `LegalBanner` est monté dans `AppLayout`, donc affiché en permanence
sur **toutes** les pages :

> Aucune donnée de ce service n'est utilisée pour déterminer vos droits ni le montant de vos allocations.

## Miroir du dépôt

Le projet est répliqué (historique complet, toutes les branches et tags) vers un
dépôt public. Depuis un clone à jour :

```bash
git push --mirror git@github.com:leolcde/CompetencePlus.git
```

`--mirror` écrase le dépôt cible avec l'état exact de la source (les branches et
tags absents en local y sont supprimés). À rejouer après chaque série de commits
à publier.
