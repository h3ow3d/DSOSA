# Dependency Policy

## Purpose

Secure Delivery Compass (DSOSA) is a security assessment tool. Its supply-chain
posture must reflect the same standards it measures. This document defines the
rules that govern every external dependency used by the project, whether at
application runtime, during the build, or in CI/release tooling.

The starting position is intentional minimal dependency. Dependencies are a
liability: each one expands the attack surface, increases review burden, and adds
operational risk. A dependency must therefore earn its place.

---

## Runtime dependency policy

A *runtime dependency* is any module, package, or asset loaded or executed when
the application is running and serving requests.

Rules:

1. **Justify before adding.** A dependency may only be introduced when there is
   a clear, documented reason why the Go standard library (or existing internal
   code) cannot meet the requirement.

2. **One decision record per dependency.** Every runtime dependency must have a
   corresponding file in `docs/dependency-decisions/` using the template in
   `docs/dependency-decisions/TEMPLATE.md`, approved before the dependency is
   added.

3. **Pin versions exactly.** Versions must be specified exactly (not as ranges or
   `latest`). Go module versions are pinned via `go.mod` and `go.sum`.

4. **Commit `go.sum`.** The `go.sum` file must always be committed so the module
   graph is reproducible and auditable.

5. **SCA coverage required.** Once SCA tooling is added to CI, every runtime
   dependency must be covered by that scan. A dependency may not bypass scanning.

6. **Reviewed updates only.** Dependency version updates must be made
   intentionally and reviewed as part of the normal PR process. Automated
   dependency update PRs must include a human review step before merging.

7. **Vulnerability response.** High or critical vulnerabilities in a runtime
   dependency must fail CI unless a written exception has been recorded (see
   Exception process below). Accepted exceptions must include a remediation
   timeline.

8. **No CDNs or remote frontend assets.** CSS, JavaScript, web fonts, and other
   static assets must be served locally from the repository. Loading assets from
   `cdn.jsdelivr.net`, `unpkg.com`, `cdnjs.cloudflare.com`,
   `fonts.googleapis.com`, or any other remote origin is forbidden.

9. **No frontend frameworks by default.** Bootstrap, Tailwind, React, Vue,
   Angular, and similar frameworks may not be added without an approved decision
   record and explicit sign-off from the project maintainers.

10. **No casual additions.** A dependency may not be added purely for convenience
    or to save a small amount of code. The convenience benefit must substantially
    outweigh the supply-chain cost.

---

## Build/test/release tooling policy

*Tooling dependencies* are tools used during CI, the build, or the release
pipeline but which do not execute in the application runtime (e.g., linters,
scanners, signing tools, container builders).

Rules:

1. Tooling used in CI must be pinned to a specific version (e.g., a full SHA for
   GitHub Actions, a tagged version for Docker images or binaries).

2. Tooling must not introduce packages into `go.mod` / `go.sum` unless they are
   also needed at application runtime.

3. New CI tooling must be documented in a commit message or PR description
   explaining its purpose.

4. Tooling that handles build artefacts or secrets must be sourced from trusted
   publishers and version-pinned.

---

## Allowed dependency categories

The following categories are considered for approval, subject to the decision
record and review process:

- Go standard library (pre-approved; no record required)
- Go modules that fill a well-defined gap not covered by the standard library
- Container base images used only in the Dockerfile (not application runtime
  packages)
- CI Actions pinned by SHA from the official GitHub Actions organisation or
  well-known publishers

---

## Forbidden dependency categories

The following are forbidden without an explicit, maintainer-approved exception:

- CDN-hosted CSS, JavaScript, or web font assets
- Frontend frameworks (Bootstrap, Tailwind CSS, React, Vue, Angular,
  govuk-frontend, Font Awesome, etc.)
- Packages added solely for convenience without a documented gap
- Packages with known high/critical unpatched CVEs
- Packages with no active maintainers or abandoned upstream
- Packages with unknown or incompatible licences

---

## Dependency approval process

1. Raise a proposal (issue or PR description) explaining:
   - What the dependency is.
   - What problem it solves.
   - Why existing code cannot solve it.
   - Alternatives considered and why they were rejected.
   - Licence, maintenance status, and transitive dependency count.

2. Complete `docs/dependency-decisions/TEMPLATE.md` and save it as
   `docs/dependency-decisions/NNNN-<name>.md`.

3. Have the decision record reviewed and approved by at least one project
   maintainer.

4. Add the dependency and verify that CI passes (including any SCA scans).

---

## Vulnerability handling

- CI must include SCA scanning once tooling is introduced (see `SUPPLY_CHAIN.md`
  for the roadmap).
- High and critical vulnerabilities discovered by SCA must block the build unless
  an exception is recorded.
- Exceptions must include: CVE identifier, affected version, impact assessment,
  and a remediation deadline.
- Medium vulnerabilities must be tracked in the issue tracker and resolved within
  a reasonable time.

---

## Update process

1. Propose the update with a brief rationale (security fix, feature need, or
   routine maintenance).
2. Review the changelog and any associated CVEs.
3. Update `go.mod` and regenerate `go.sum`.
4. Ensure all tests pass.
5. Merge through the normal PR review process.

---

## Exception process

An exception is required when:

- A dependency is added that does not meet all criteria above.
- A vulnerability is accepted rather than immediately remediated.

To record an exception:

1. Create `docs/dependency-decisions/exceptions/NNNN-<description>.md`.
2. Include: what is being excepted, the risk accepted, the justification, the
   expiry date or remediation target, and the approver.
3. Get explicit sign-off from a project maintainer.

Exceptions are time-limited and must be revisited at or before their expiry date.
