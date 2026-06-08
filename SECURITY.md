# Security Policy

## Scope

Agent-Tools is a **certifying registry**. It does not publish, warrant, or
verify the third-party tools it catalogs — that is the publishing platform
owner's responsibility (see
[`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md#trust--certification-model)).
This policy covers the code and automation **in this repository**.

## Principles

- **Few trusted platforms, maximized.** Our security supply chain is essentially
  one vendor — **GitHub**, our trusted platform owner — used to its fullest.
- **One task, one tool. No overlapping responsibility.** Each security concern
  has exactly one owner.

## Responsibility map (one owner per task)

| Task | Owner | Vendor |
|---|---|---|
| SAST — vulnerabilities in our code (Go + Python) | **CodeQL** (default setup) | GitHub |
| Dependency vulnerabilities + remediation | **Dependabot** (alerts, security & version updates) | GitHub |
| Dependency / Actions version freshness | **Dependabot** (`.github/dependabot.yml`) | GitHub |
| Secret detection | **Secret scanning + push protection** | GitHub |
| Build integrity / provenance | **`actions/attest-build-provenance`** (`provenance.yml`) | GitHub |
| Build, test & registry-rule enforcement | **`ci.yml`** (`go test` + `validate_registry.py`) | self |

Nothing else is added: extra scanners (gosec, Semgrep, Trivy, Gitleaks,
govulncheck, dependency-review, Scorecard) were intentionally **left out** —
each duplicated a task already owned above and/or introduced another vendor.

## Maximize GitHub: settings to enable (platform-owner toggles)

These native features have no in-repo file; enable them in **Settings → Code
security**:

- **Secret scanning** + **push protection**
- **Dependabot security updates** (auto-fix PRs for vulnerable deps)
- **CodeQL**: add the **`actions`** language so workflows are scanned too
- **Private vulnerability reporting**
- A **ruleset on `main`** requiring the `ci.yml` check, review, and linear history

## Build provenance

`provenance.yml` produces a signed SLSA build-provenance attestation for the
`worker` binary — GitHub (the platform owner) verifying our build, exactly as
the trust model prescribes. Verify with:

```bash
gh attestation verify dist/worker --repo AGenNext/Agent-Tools
```

## Reporting a vulnerability

Report privately via **GitHub → Security → Report a vulnerability** rather than a
public issue. We aim to acknowledge within a few business days.
