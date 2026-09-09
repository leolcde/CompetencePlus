# ProfilsActifs

Monorepo : `backend/` (API Go), `frontend/` (Vue 3 + Vite), `database/` (migrations SQL).

## Prérequis

- Docker + Docker Compose
- Pour un lancement manuel : Go 1.25+, Node 22+, PostgreSQL 16

## Configuration

Copier le fichier d'exemple et le remplir :

```bash
cp .env.example .env
```

| Variable            | Description                                   |
| ------------------- | --------------------------------------------- |
| `POSTGRES_USER`     | utilisateur PostgreSQL                         |
| `POSTGRES_PASSWORD` | mot de passe PostgreSQL                        |
| `POSTGRES_DB`       | nom de la base                                 |
| `DATABASE_URL`      | DSN complet utilisé par le backend            |
| `JWT_SECRET`        | secret de signature des tokens JWT            |
| `PORT`              | port d'écoute du backend (def. `8080`)        |

> `DATABASE_URL` doit rester cohérent avec les identifiants `POSTGRES_*`.

## Lancer avec Docker (recommandé)

```bash
docker compose up --build
```

- Frontend : http://localhost:5173
- Backend : http://localhost:8080 (health : `GET /health`)
- PostgreSQL : `localhost:5432`

Compose n'applique pas les migrations automatiquement — voir la section Migrations.

Arrêter : `docker compose down` (ajouter `-v` pour supprimer le volume de données).

## Lancer manuellement

### 1. Base de données

Démarrer un PostgreSQL local, puis appliquer les migrations avec
[golang-migrate](https://github.com/golang-migrate/migrate) :

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path database/migrations -database "$DATABASE_URL" up
```

La migration `000006_seed_fake_profiles` insère des comptes de démo.

### 2. Backend

```bash
cd backend
go mod download
go run .            # ou : go build -o dist/profilsactifs-server . && ./dist/profilsactifs-server
```

Endpoints : `/`, `/health`, `/auth/register`, `/auth/login`.

### 3. Frontend

```bash
cd frontend
npm install
npm run dev         # http://localhost:5173
```

## Build de production

```bash
# Backend
cd backend && go build -o dist/profilsactifs-server .

# Frontend
cd frontend && npm run build     # sortie dans frontend/dist
```

## Migrations

```bash
migrate -path database/migrations -database "$DATABASE_URL" up      # appliquer
migrate -path database/migrations -database "$DATABASE_URL" down 1  # revenir en arrière
```

## Tests

```bash
cd backend && go vet ./... && go test ./...
```

## CI

`.github/workflows/build.yml` s'exécute sur push `dev` et PR vers `_access_test` :
démarre un service PostgreSQL, applique les migrations, lance `go vet` / `go test`,
compile le binaire et l'upload en artifact.