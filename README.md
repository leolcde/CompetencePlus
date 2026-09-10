<h1 align="center">ProfilsActifs</h1>

<p align="center">
  <b>Plateforme de mise en relation entre demandeurs d'emploi et recruteurs</b><br/>
  Présentation vidéo · questionnaire de savoir-être · badge de certification
</p>

<p align="center">
  <img src="https://img.shields.io/badge/backend-Go%201.26-00ADD8?style=flat-square&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/frontend-Vue%203%20%2B%20Vite-4FC08D?style=flat-square&logo=vuedotjs&logoColor=white" />
  <img src="https://img.shields.io/badge/db-PostgreSQL%2016-4169E1?style=flat-square&logo=postgresql&logoColor=white" />
  <img src="https://img.shields.io/badge/run-Docker%20Compose-2496ED?style=flat-square&logo=docker&logoColor=white" />
</p>

---

## À propos

**ProfilsActifs** permet aux demandeurs d'emploi de se présenter en **vidéo** (au-delà du CV) et
aux recruteurs de parcourir ces profils. Un **questionnaire de savoir-être** délivre un
**badge de certification**. Le consentement à la diffusion de la vidéo est explicite et révocable
(droit à l'effacement).

## Lancer le projet

```bash
cp .env.example .env      # renseigner POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_DB / JWT_SECRET
docker compose up --build
```

| Service    | URL                            |
| ---------- | ------------------------------ |
| Interface  | http://localhost:5173          |
| API        | http://localhost:8080          |
| État API   | http://localhost:8080/health   |

Arrêt : `docker compose down` (ajouter `-v` pour repartir d'une base vierge).

La base est créée automatiquement au premier démarrage, le questionnaire est pré-rempli.
Il suffit ensuite de **créer un compte** depuis l'interface.

## Utilisation

1. **Inscription** — nom, e-mail, mot de passe, date de naissance (16 ans minimum).
2. **Questionnaire** — 20 questions de savoir-être → un badge si le score est suffisant.
3. **Ma vidéo** — coller un lien (YouTube, Vimeo…) ou téléverser un fichier, après avoir donné son consentement.
4. **Candidats** — la liste et les fiches sont consultables publiquement.

## Documentation

| Document | Contenu |
| --- | --- |
| [devdoc.md](devdoc.md) | architecture, API, variables d'env, lancement manuel, build |
| [SECURITE.md](SECURITE.md) | mesures de sécurité et limitations connues |
| [TEST.md](TEST.md) | procédure de test |
