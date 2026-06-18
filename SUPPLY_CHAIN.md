# Supply Chain Trust Model

## Purpose

This document describes how Secure Delivery Compass (DSOSA) is built, delivered,
and verified. It defines the trust boundaries, current controls, and the roadmap
for future supply-chain hardening.

---

## Source of truth

The canonical source for all application code, configuration, documentation, and
CI definitions is this Git repository. No code is authoritative unless it has
been committed here and has passed CI.

- All changes enter through pull requests.
- Direct pushes to the main branch are discouraged and should be protected by
  branch protection rules.
- Commit history must not be rewritten on the main branch.

---

## CI as the controlled build environment

All build, test, and verification steps run inside GitHub Actions. Local builds
are used only for development; they must not produce artefacts that are published
or deployed.

CI is responsible for:

- Compiling the application (`go build`).
- Running the test suite (`go test ./...`).
- Enforcing code formatting (`gofmt`).
- Running static analysis (`go vet`).
- Building the Docker image.
- Running the smoke test.
- Enforcing the dependency baseline (see below).

CI tooling (Actions, Docker, scanning tools) may use external tools. Those tools
are considered *CI tooling*, not application runtime dependencies. They must be
pinned to specific versions and must not introduce packages into the application's
`go.mod`.

---

## Dependency policy

All dependency decisions are governed by [DEPENDENCY_POLICY.md](DEPENDENCY_POLICY.md).

Key principles:

- External dependencies are allowed only when intentional, minimal, pinned,
  auditable, and scanned.
- Every application runtime dependency requires an approved decision record in
  `docs/dependency-decisions/`.
- CI/release tooling may use external tools, but those tools must not become
  accidental runtime dependencies.
- No CDNs or remote frontend assets are permitted in the application.

---

## Current runtime artefact expectations

The production artefact is a single Docker image built from `Dockerfile`.

Expected content:

| Component          | Origin                      | Notes                                      |
|--------------------|-----------------------------|--------------------------------------------|
| Go binary          | Built from source in CI     | Statically linked; no external Go modules  |
| HTML templates     | `web/templates/`            | Server-rendered; no external frameworks    |
| CSS / JS assets    | `web/static/`               | Local only; no CDN references              |
| Base image         | `golang:1.23-alpine` (build)| Multi-stage build; final image is scratch/alpine |
| Runtime data       | Mounted volume (`/data`)    | Not baked into the image                   |

The application currently has **no third-party Go module dependencies**. This is
intentional and documented in
`docs/dependency-decisions/0000-current-baseline.md`.

---

## Sensitivity of /data

The `/data` mount contains assessment results, project records, and the event
log. This data:

- Must not be baked into container images.
- Must be excluded from version control (see `.gitignore`).
- Should be encrypted at rest in production environments.
- Should be backed up regularly in production environments.
- May contain sensitive organisational maturity information and must be treated
  accordingly.

---

## Planned supply-chain controls

The following controls are planned and will be added in future steps. They are
listed here to make the intended security posture explicit.

### Secret scanning

Automated scanning of commits and pull requests for accidentally committed
secrets (tokens, keys, credentials). Will use a GitHub Advanced Security feature
or an equivalent open-source tool.

### SAST (Static Application Security Testing)

Static analysis of application source code for security vulnerabilities. Planned
tool: CodeQL (via GitHub Actions) or equivalent.

### SCA (Software Composition Analysis)

Scanning of application dependencies for known CVEs. Will integrate with
`go.sum` and the Go vulnerability database (`govulncheck`). High/critical
findings will fail CI.

### DAST (Dynamic Application Security Testing)

Automated scanning of the running application for common web vulnerabilities.
Will be added once the application has a stable test environment.

### IaC scanning

Scanning of `Dockerfile` and `docker-compose.yml` for misconfigurations. Planned
tool: Hadolint (Dockerfile) and Checkov or equivalent.

### Container scanning

Scanning of the built Docker image for known CVEs in OS packages and language
runtimes. Planned tool: Trivy or equivalent, run in CI after `docker build`.

### SBOM (Software Bill of Materials)

Generation of an SBOM for each release artefact, documenting all components and
their versions.

### Signing

Cryptographic signing of release artefacts (container images and/or binaries)
using Sigstore/cosign or equivalent.

### Attestations

Build provenance attestations (SLSA level 2 or higher) linking artefacts back to
the specific source commit and CI run that produced them.

---

## Separation between runtime and CI tooling

The project distinguishes two categories of external tool:

| Category               | Examples                          | Location           |
|------------------------|-----------------------------------|--------------------|
| App runtime dependency | Go modules in `go.mod`            | Ships in the image |
| CI/release tooling     | GitHub Actions, scanners, linters | CI only; not shipped |

CI tooling must not appear in `go.mod` unless it is also a genuine application
runtime dependency. Scanning tools, linters, and release utilities are CI
tooling only.

---

## Future release verification approach

When artefact signing and attestation are implemented, the intended verification
approach is:

1. Download the release artefact (container image or binary).
2. Verify the cryptographic signature against the project's published public key.
3. Verify the build provenance attestation against the source commit SHA.
4. Optionally verify the SBOM against known-good component hashes.

This approach allows consumers to verify that an artefact was built from a
specific, reviewed commit in this repository by the project's CI pipeline, and
not tampered with afterwards.
