# Dependency Decision: Current MVP Baseline

## Dependency name

(none — this record documents the intentional absence of third-party application
dependencies at the start of the project)

## Version

n/a

## Type

Baseline record

## Purpose

To document that the Secure Delivery Compass MVP is intentionally built using
only the Go standard library, local assets, and Docker/GitHub Actions as delivery
tooling, with no third-party application runtime dependencies.

## Why it is needed

This record establishes the starting point of the dependency decision log. It
makes the current state explicit so that any future addition of a runtime
dependency is clearly visible as a deliberate change from this baseline.

## Alternatives considered

Numerous Go libraries exist for web frameworks, template engines, JSON handling,
and storage. All were evaluated against the question: "Can the Go standard library
do this?"

## Why alternatives were rejected

At MVP stage, every requirement — HTTP routing, HTML templating, JSON
marshalling/unmarshalling, file I/O, and SHA hashing — is satisfied by the Go
standard library. Introducing external modules at this stage would add supply-chain
risk without delivering a corresponding benefit.

## Runtime impact

There are no third-party runtime dependencies. The production binary links only
against the Go standard library. This results in a minimal attack surface and a
fully reproducible build.

## Security considerations

No third-party code executes at application runtime, which eliminates the most
common vector for supply-chain attacks (compromised upstream packages). The Go
standard library is maintained by the Go team and covered by the Go security
disclosure process.

## Licence considerations

The Go standard library is covered by the BSD-style Go licence, which is
compatible with all common open-source licensing scenarios.

## Maintenance status

The Go standard library is actively maintained by Google and the Go open-source
community with a documented backwards-compatibility guarantee.

## Transitive dependencies

None. `go list -m all` returns only the root module:

```
dsovs-assessment-tool
```

## CI/SCA coverage

SCA scanning is planned (see `SUPPLY_CHAIN.md`). Because there are currently no
third-party modules, SCA will return no findings. When SCA is enabled, it will
cover any future modules introduced.

## Application component inventory

The following components make up the application at this baseline:

| Component                  | Origin                             | Notes                               |
|----------------------------|------------------------------------|-------------------------------------|
| Go application code        | This repository (`cmd/`, `internal/`) | Standard library only               |
| HTML templates             | `web/templates/`                   | Server-rendered; no JS frameworks   |
| CSS                        | `web/static/styles.css`            | Local; no CDN; no external fonts    |
| JavaScript                 | `web/static/app.js`                | Local; minimal; no frameworks       |
| File-backed storage        | `internal/storage/`                | JSON files on local filesystem      |
| Docker (delivery tooling)  | `Dockerfile`, `docker-compose.yml` | CI and local dev only               |
| GitHub Actions (CI tooling)| `.github/workflows/ci.yml`         | CI only; not a runtime dependency   |

## Rollback plan

n/a — this is the baseline. Any future dependency that needs to be rolled back
will have its own decision record with a rollback plan.

## Approval date

2026-06-18
