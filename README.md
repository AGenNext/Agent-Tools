# Agent-Tools

> A **certifying tool-registry** for AI agents: container & Kubernetes
> capabilities (podman, kind, k3s) exposed as orchestrated tools, plus a
> complete catalog of the SDKs, CLIs, and MCP servers for the orchestration and
> developer-portal platforms around them (Orkes Conductor, Temporal, Cortex).

**Keywords:** agent tools · tool registry · MCP servers · Orkes Conductor ·
Temporal · Cortex IDP · podman · kind · k3s · Kubernetes · workflow
orchestration · durable execution · GitOps catalog · software supply-chain
certification.

> 🤖 **Agents:** read [`AGENTS.md`](AGENTS.md) first — it explains the repo and
> the rules so you don't need to scan the whole tree.

## What this is

Agent-Tools turns tools into **registered, orchestrated, catalogued, and
certified** entries. Three layers (see [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md)):

| Layer | What | Tech |
|---|---|---|
| Capabilities | binary-backed tools an agent runs | podman · kind · k3s (Go workers) |
| Orchestration / execution | durable, multi-step workflows | Orkes Conductor · Temporal (both first-class) |
| Operational excellence | catalog, scorecards, governance | Cortex (IDP) via GitOps |

It is a **certifier, not a publisher**: it records each tool's version and
capabilities *as claimed by the publisher*, trusting the **platform owner**
(e.g. GitHub) as the verifier. See the
[trust model](docs/ARCHITECTURE.md#trust--certification-model).

## Rules (invariants, enforced in CI — see [`docs/TAXONOMY.md`](docs/TAXONOMY.md))

1. **No-blocker:** every tool and platform declares ≥1 alternative — nothing is a single point of failure.
2. **Certification, not publication:** claims are the publisher's; verification is the platform owner's duty.
3. **No fabrication:** versions are real published releases; changelogs point at exactly that release.
4. **Controlled vocabulary** for `type`/`provenance`; **editorial fields labelled**.

## Registry contents

- **Capabilities (3):** `podman`, `kind`, `k3s` — each a Conductor worker with a
  [skill](skills/) and a Cortex [Tool-Card](skills/podman/cortex.yaml).
- **Catalog (19 tools across 3 apps):** every-language SDK, CLI, and MCP server
  for Conductor / Temporal / Cortex — [`catalog/apps.yaml`](catalog/apps.yaml),
  summarized in [`docs/CATALOG.md`](docs/CATALOG.md).

## Layout

```
internal/tools/        # capability workers + Registry (source of truth) + invariants test
cmd/worker/            # Orkes Conductor task runner that registers the workers
skills/<tool>/         # SKILL.md (runbook) + cortex.yaml (Tool-Card)
catalog/apps.yaml      # catalog of all SDKs/CLIs/MCPs for the apps
templates/             # tool-card template
scripts/               # validate_registry.py (rule enforcement)
docs/                  # ARCHITECTURE, TAXONOMY, CATALOG
.github/               # CI, dependabot, issue/PR templates, CODEOWNERS
```

## Run the orchestrator worker

```bash
export CONDUCTOR_SERVER_URL=https://play.orkes.io/api
export CONDUCTOR_AUTH_KEY=...      # Orkes → Access Control → Applications
export CONDUCTOR_AUTH_SECRET=...
go run ./cmd/worker
```

## Develop

```bash
go build ./... && go vet ./... && go test ./...   # Go + registry invariants
python scripts/validate_registry.py               # taxonomy + no-blocker rules
```

CI runs the same checks and fails on any rule violation. Contributions:
[`CONTRIBUTING.md`](CONTRIBUTING.md).
