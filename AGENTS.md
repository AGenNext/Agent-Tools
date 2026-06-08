# AGENTS.md

Orientation for AI agents (and humans) working in this repo. Read this and the
linked docs **instead of scanning the whole tree**.

## What this repo is

Agent-Tools is a **certifying tool-registry**. It exposes container/Kubernetes
capabilities (podman, kind, k3s) as orchestrated tools and catalogs the SDKs,
CLIs, and MCP servers of the platforms around them (Orkes Conductor, Temporal,
Cortex). It **certifies** entries — it does not publish or verify the tools
themselves.

Full picture: [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) (layers + trust
model), [`docs/TAXONOMY.md`](docs/TAXONOMY.md) (vocabulary + rules),
[`docs/CATALOG.md`](docs/CATALOG.md) (catalogued tools).

## Rules you must respect (enforced by CI — do not break)

1. **No-blocker:** every tool/platform declares ≥1 `alternatives`. Nothing is a blocker.
2. **Certification, not publication:** record `version`/`capabilities` as *publisher-claimed*; the **platform owner** (e.g. GitHub) is the verifier. Never claim you verified a tool.
3. **No fabrication:** pin real versions (`git ls-remote`); the changelog must point at exactly that release. Omit unknowns.
4. **Controlled vocabulary:** `type` ∈ {capability, sdk, cli, mcp, gitops, api}; `provenance` starts with official/community.
5. **Label editorial fields** (`alternatives`, opinions) as Agent-Tools editorial.
6. **Keep the path optimal & living:** update the roadmap in `docs/ARCHITECTURE.md` as new information arrives.

## Where things live

| Need | Go here |
|---|---|
| Capability worker logic | `internal/tools/` (`tools.go` Registry = source of truth; `exec.go` runner; `runner.go` Conductor adapter) |
| Orchestrator entrypoint | `cmd/worker/main.go` |
| Skill / runbook per tool | `skills/<tool>/SKILL.md` |
| Cortex Tool-Card per tool | `skills/<tool>/cortex.yaml` (template in `templates/`) |
| Catalog of SDKs/CLIs/MCPs | `catalog/apps.yaml` |
| Rule enforcement | `scripts/validate_registry.py` + `internal/tools/registry_test.go` |

## How to add a tool (short form)

1. `internal/tools/tools.go`: add a `Registry` entry (real version, publisher,
   verified-by, changelog, capabilities, **alternatives ≥1**, worker) + worker func.
2. `skills/<tool>/SKILL.md` and `skills/<tool>/cortex.yaml` (from the template).
3. Validate, then open a PR (template + checklist apply).

Catalog-only entries (SDK/CLI/MCP) go in `catalog/apps.yaml`. Details:
[`CONTRIBUTING.md`](CONTRIBUTING.md).

## Always validate before pushing

```bash
go build ./... && go vet ./... && go test ./...
python scripts/validate_registry.py
```

Both must pass — CI runs them and blocks rule violations.
