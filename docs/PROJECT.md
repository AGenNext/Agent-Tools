# Project Tracker

We build in **loops** — small, non-overlapping increments — not the whole
catalog at once. This is the living status board.

## Done (PR #2)

- [x] **Capabilities:** podman / kind / k3s as Conductor workers (`internal/tools`, `cmd/worker`); skills + runbooks
- [x] **Catalog:** all SDKs / CLIs / MCP servers for Conductor, Temporal, Cortex (`catalog/apps.yaml`, `docs/CATALOG.md`)
- [x] **Trust model:** publisher claims → platform owner verifies → Agent-Tools certifies; integration mode decides who validates; prefer MCP/API over SDK (`docs/ARCHITECTURE.md`)
- [x] **Taxonomy & rules:** controlled vocabulary + invariants — no-blocker, no fabrication, certification (`docs/TAXONOMY.md`)
- [x] **Real versions** pinned via `git ls-remote`; Tool-Cards with publisher / verified-by / changelog / alternatives
- [x] **Single source of truth:** Go Registry is canonical; a test enforces the Tool-Cards mirror it (`internal/tools/cards_test.go`)
- [x] **Enforcement:** `scripts/validate_registry.py` + `go test` invariants, wired into `ci.yml`
- [x] **Governance:** CONTRIBUTING, issue/PR templates, CODEOWNERS, card template
- [x] **Security (consolidated, one task/one tool):** CodeQL, Dependabot, secret scanning, build provenance (`provenance.yml`); `SECURITY.md`
- [x] **Discoverability:** README + AGENTS.md; keyword-rich
- [x] **Dep hygiene:** bumped `protobuf` → v1.34.2 (moderate CVE)
- [x] PR #2 marked ready for review

## Blocked — needs platform/org owner (not code)

- [ ] **Org Actions policy** blocks self-defined workflows (`ci.yml`, `provenance.yml` → `startup_failure`). Fix: *Settings → Actions → Allowed actions → Allow all* (CodeQL/Dependabot run regardless).
- [ ] **Repo About:** set description + topics (no MCP tool for it; values in README).
- [ ] **Code-security toggles:** secret scanning push protection, Dependabot security updates, CodeQL `actions` language, branch ruleset requiring `ci.yml` (see `SECURITY.md`).

## Next loops (one each — pick per priority)

- [ ] **Conductor task-definition registration** — a `register` command using the Conductor SDK so task/workflow defs exist server-side (extends an in-supply-chain platform; no overlap).
- [ ] **Worked example** — add one new tool end-to-end (registry entry + skill + card) as the canonical contribution demo.
- [ ] **Run-through** — stand up Conductor (play.orkes.io) and execute one tool to prove the worker end-to-end.

## Decisions parked

- **Temporal stays catalog-only.** Embedding its SDK would add a second,
  overlapping workflow engine to our supply chain — revisit only if a concrete
  need requires code-first long-running sagas Conductor can't serve.
