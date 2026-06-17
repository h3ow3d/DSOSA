# DSOVS Assessment Tool

A local-first **DevSecOps Maturity Assessment** web app built on the
[OWASP DevSecOps Verification Standard (DSOVS)](https://owasp.org/www-project-devsecops-verification-standard/).

---

## Architecture

| Layer | Technology |
|-------|-----------|
| Frontend | React 18 · TypeScript · Vite · Tailwind CSS · Recharts |
| Backend | Python FastAPI · SQLAlchemy 2 · Alembic · Pydantic v2 |
| Database | PostgreSQL 16 |
| Orchestration | Docker Compose |

```
dsovs-assessment-tool/
  docker-compose.yml
  .env.example
  backend/   – FastAPI app, SQLAlchemy models, Alembic migrations
  frontend/  – React SPA
```

---

## How to Run Locally

### Prerequisites

* Docker Desktop (or Docker Engine + Compose v2)

### 1. Clone and configure

```bash
git clone <repo-url>
cd dsovs-assessment-tool
cp .env.example .env
# Edit .env if needed (defaults work out of the box)
```

### 2. Start the stack

```bash
docker compose up --build
```

* Frontend: <http://localhost:5173>
* Backend API docs: <http://localhost:8000/docs>

---

## How to Sync OWASP DSOVS Data

1. Open the app at <http://localhost:5173>
2. Click **"Sync DSOVS Catalogue"** on the Dashboard
3. The app will fetch the latest DSOVS JSON from OWASP, calculate a SHA256
   hash, and store it in the database

The sync endpoint is idempotent – re-running with the same data is a no-op.

To auto-sync on backend startup set `AUTO_SYNC_CATALOGUE=true` in `.env`.

---

## How to Create a Project

1. Navigate to **Projects** → click **+ New Project**
2. Fill in the project name and optional metadata
3. Click **Create Project**

---

## How to Create an Assessment

1. Open a project from the Projects list
2. Click **+ New Assessment**
3. Fill in the assessment name, assessor, and scope
4. Click **Create & Start** – you will be taken to the Assessment Wizard

---

## Assessment Wizard

The wizard groups all DSOVS controls by phase. For each control you can:

* Select a current maturity level (0–3)
* Select a target maturity level
* Mark a control as Not Applicable
* Add evidence and action notes
* Set priority and confidence

Scores are **auto-saved** with a 600 ms debounce.

---

## How to Print / Save a PDF Report

1. Complete an assessment
2. Click **Results** to open the results page
3. Click **Print / Save PDF** – this opens a new tab with the printable report
4. Use your browser's **File → Print** or **Ctrl+P** / **Cmd+P**
5. Select **Save as PDF** as the destination

---

## API Route Summary

### Catalogue

| Method | Route | Description |
|--------|-------|-------------|
| `POST` | `/api/catalogue/sync` | Fetch & store DSOVS catalogue |
| `GET`  | `/api/catalogue/current` | Return latest catalogue |

### Projects

| Method | Route | Description |
|--------|-------|-------------|
| `GET`    | `/api/projects` | List projects |
| `POST`   | `/api/projects` | Create project |
| `GET`    | `/api/projects/{id}` | Get project |
| `PUT`    | `/api/projects/{id}` | Update project |
| `DELETE` | `/api/projects/{id}` | Delete project |

### Assessments

| Method | Route | Description |
|--------|-------|-------------|
| `GET`    | `/api/projects/{id}/assessments` | List assessments |
| `POST`   | `/api/projects/{id}/assessments` | Create assessment |
| `GET`    | `/api/assessments/{id}` | Get assessment |
| `PUT`    | `/api/assessments/{id}` | Update assessment |
| `DELETE` | `/api/assessments/{id}` | Delete assessment |

### Scores & Results

| Method | Route | Description |
|--------|-------|-------------|
| `PUT` | `/api/assessments/{id}/scores/{control_id}` | Upsert score |
| `GET` | `/api/assessments/{id}/results` | Get scored results |
| `GET` | `/api/projects/{id}/trends` | Get maturity trend |
| `GET` | `/api/assessments/{id}/report-data` | Get full report data |

---

## Database Migration Commands

```bash
# Inside the backend container
docker compose exec backend alembic upgrade head

# Create a new migration after changing models
docker compose exec backend alembic revision --autogenerate -m "description"
```

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| Frontend shows "No DSOVS catalogue loaded" | Click **Sync DSOVS Catalogue** on the Dashboard |
| Cannot create assessment | Sync the catalogue first |
| Backend fails to start | Ensure `db` container is healthy; check `DATABASE_URL` in `.env` |
| CORS errors in browser | Ensure `FRONTEND_ORIGIN` in `.env` matches the frontend URL |
| Port conflict | Change ports in `docker-compose.yml` or `.env` |

---

## Environment Variables

See `.env.example` for all supported variables.

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgresql+psycopg://dsovs:dsovs@db:5432/dsovs` | PostgreSQL connection string |
| `DSOVS_API_URL` | OWASP DSOVS JSON URL | Source for catalogue sync |
| `FRONTEND_ORIGIN` | `http://localhost:5173` | CORS allowed origin |
| `AUTO_SYNC_CATALOGUE` | `false` | Sync catalogue on startup |
| `VITE_API_BASE_URL` | `http://localhost:8000` | Frontend → backend base URL |
