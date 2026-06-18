# DSOVS Assessment Tool (Go Native)

Local-first OWASP DSOVS maturity assessment tool implemented as a single Go web application using only the Go standard library.

## Architecture

```text
  docker-compose.yml
  Dockerfile
  README.md
  .env.example
  cmd/
    dsovs/
      main.go
  internal/
    app/
      server.go
      routes.go
    dsovs/
      client.go
      models.go
      sync.go
    storage/
      store.go
      json_store.go
      events.go
    assessment/
      models.go
      scoring.go
      reports.go
      improvements.go
    web/
      templates.go
      static.go
  web/
    templates/
      layout.html
      dashboard.html
      projects.html
      project.html
      assessment.html
      results.html
      report.html
    static/
      styles.css
      app.js
```

## Runtime model

- Single service in Docker Compose.
- App served at `http://localhost:8080`.
- File persistence mounted under `/data`.
- Catalogue sync source:
  `https://owasp.org/www-project-devsecops-verification-standard/dist/dsovs.json`

Persisted data layout:

- `/data/catalogue/` (versioned catalogue snapshots)
- `/data/projects/`
- `/data/assessments/`
- `/data/improvements/`
- `/data/events.ndjson` (append-only event log)

## Run

```bash
cp .env.example .env
docker compose up --build
```

Then open `http://localhost:8080`.

## Testing the UI

- Run all tests: `go test ./...`
- Run the app locally and open `http://localhost:8080`
- Manually test keyboard navigation (including skip link and focus order)
- Manually test browser print and Save as PDF from a report page
- Manually verify pages remain readable and usable with JavaScript disabled

## Dependency governance

Secure Delivery Compass allows dependencies only through a governed process.
Application dependencies must be documented, pinned, reviewed, scanned, and
covered by CI guardrails. CI/release tooling may use external tools, but those
tools must not become accidental runtime dependencies.

- [DEPENDENCY_POLICY.md](DEPENDENCY_POLICY.md) — rules for every external dependency
- [SUPPLY_CHAIN.md](SUPPLY_CHAIN.md) — delivery trust model and planned controls
- [docs/dependency-decisions/](docs/dependency-decisions/) — per-dependency decision records

## Dashboard sync flow

1. Open Dashboard.
2. Click **Sync OWASP DSOVS Catalogue**.
3. The app fetches DSOVS JSON, computes SHA256, stores a new versioned file in `/data/catalogue/`, and appends a sync event to `/data/events.ndjson`.
4. Re-sync with unchanged data is a no-op.
