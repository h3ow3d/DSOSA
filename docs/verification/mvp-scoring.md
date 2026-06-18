# MVP Scoring Journey — Verification Checklist

This document records the expected behaviour of the MVP scoring journey.
All items must pass before any change is merged to the main branch.

## Manual verification steps

### 1. Catalogue sync

- Start the application (`docker compose up --build` or `go run ./cmd/dsovs`).
- Open `http://localhost:8080`.
- Click **Sync OWASP DSOVS Catalogue**.
- Verify that the sync succeeds and the catalogue version is shown on the
  dashboard.
- Re-sync with unchanged data and verify it is a no-op (no duplicate entries).

### 2. Project creation

- Navigate to the Projects page.
- Create a new project with a name and description.
- Verify the project appears in the project list.

### 3. Assessment creation

- Open a project.
- Create a new assessment.
- Verify the assessment appears in the assessment list for that project.

### 4. Scoring — ORG-001

- Open the assessment.
- Locate control **ORG-001** (or the first available control in the catalogue).
- Select a score (e.g., **Implemented**).
- Save the score.
- Verify that the score is saved without errors.

### 5. Saved score persists after refresh

- After saving ORG-001, reload the assessment page.
- Verify that the previously saved score is still shown correctly.

### 6. Results page reflects saved scores

- Navigate to the Results page for the assessment.
- Verify that the saved score for ORG-001 (and any other scored controls) is
  reflected in the results summary.

## Automated verification

Run from the repository root:

```bash
go test ./...
```

All tests must pass.

## Docker Compose

```bash
docker compose up --build
```

The application must start successfully and be reachable at
`http://localhost:8080`.
