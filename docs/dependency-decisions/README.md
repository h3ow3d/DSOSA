# Dependency Decisions

This directory records every intentional application dependency decision for
Secure Delivery Compass (DSOSA).

## Why this exists

Every external dependency carries supply-chain risk. Recording decisions here
ensures that:

- Each dependency is added intentionally, not casually.
- The rationale, alternatives, and trade-offs are visible to all contributors.
- Future reviewers can understand why a dependency was chosen.
- Dependency reviews and audits have a clear paper trail.

## When to add a record

A decision record is required **before** adding any of the following to the
application:

- A new entry in `go.mod` (third-party Go module).
- A new CDN-hosted asset referenced in templates or static files (currently
  forbidden; see [DEPENDENCY_POLICY.md](../../DEPENDENCY_POLICY.md)).
- A new frontend library or framework.
- Any other external component that ships in the production artefact or is loaded
  at runtime.

Decision records are **not** required for:

- The Go standard library.
- CI-only tooling that does not affect `go.mod` or the production artefact.

## How to add a record

1. Copy `TEMPLATE.md` to `NNNN-<dependency-name>.md` (e.g.,
   `0001-some-library.md`), using the next available four-digit number.
2. Fill in every section of the template.
3. Submit the record as part of the PR that introduces the dependency.
4. Get approval from at least one project maintainer before merging.

## Index

| File                          | Dependency            | Status   |
|-------------------------------|-----------------------|----------|
| 0000-current-baseline.md      | (no dependencies)     | Active   |
