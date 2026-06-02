# Agent-Tools Architecture

This repo turns container / Kubernetes capabilities (podman, kind, k3s) into
**orchestrated agent tools**. The system is layered so each technology owns a
distinct concern — and so we don't run two overlapping workflow engines.

```
┌──────────────────────────────────────────────────────────┐
│  OPERATIONAL EXCELLENCE        →  Cortex (IDP)             │
│  catalog every tool-worker as a service, ownership,       │
│  scorecards (runbook? healthy? on-call?), standards       │
├──────────────────────────────────────────────────────────┤
│  ORCHESTRATION / EXECUTION     →  Conductor                │
│  durable, multi-step workflows that chain the tools       │
├──────────────────────────────────────────────────────────┤
│  CAPABILITIES (agent tools)    →  Go workers               │
│  podman / kind / k3s — run the CLI, return stdout/exit     │
└──────────────────────────────────────────────────────────┘
```

## Why one workflow engine, not two

Orkes **Conductor** and **Temporal** occupy the *same* layer — both are durable
workflow / execution engines. Running both ("Conductor orchestrates, Temporal
executes") means two servers, two SDKs, and two failure models for marginal
benefit.

For this project each tool is a discrete, pollable capability that an agent
composes into pipelines, so **Conductor** is the chosen engine: declarative
workflows, each tool is a registered worker, trivial to add tools, full visual
run history.

**Temporal is deferred.** Add it only if a single tool operation becomes a
complex, long-running, stateful saga (retries / heartbeats / human approval
over hours–days). At that point a Conductor task can delegate to a Temporal
workflow — a bridge to cross when a concrete tool needs it, not before.

**Cortex** is *not* a duplicate of either — it never executes work. It is the
Internal Developer Portal layer: it catalogs the tool-workers as services and
runs scorecards for operational excellence.

## Trust & certification model

Agent-Tools is a **certifying registry, not a publisher**. It records each
tool's version and capabilities and certifies that the entry faithfully
references the upstream source — nothing more. Responsibilities are split three
ways:

| Role | Who | Responsibility |
|------|-----|----------------|
| **Publisher** | the tool's vendor/org (e.g. containers, kubernetes-sigs, k3s-io) | *Makes* the claims — declares versions and capabilities in its changelog. |
| **Verifier** | the **platform owner** (currently GitHub) | *Verifies* those claims. Authenticity is, everywhere, the platform owner's responsibility and duty — via signed tags, verified releases, and build provenance / attestations. |
| **Certifier** | Agent-Tools | *Trusts* the platform owner as the verification authority and certifies only that the registry entry references the platform-hosted source. Does not publish, warrant, or independently verify the tool. |

Concretely: a Tool-Card's `version` and `capabilities` are **as claimed by the
publisher**; the `changelog` is the platform-hosted evidence; and we trust the
platform owner (GitHub) to have verified it.

**Where we're headed (and hope to soon):** when Agent-Tools operates its own
platform, it *becomes* the platform owner and assumes the verification duty
itself, rather than delegating it to an external platform. Until then, we trust
the platform owner as the verifier.

## Components

### Capabilities — `internal/tools`

Each agent tool is a Conductor worker function with the signature
`func(*model.Task) (interface{}, error)`. `runner.go` provides a shared,
shell-free executor that runs a fixed binary with an `args` array (so task
input can't inject extra commands), honors an optional `timeoutSeconds`, pipes
optional `stdin`, and returns a structured `Result` (binary, args, exitCode,
stdout, stderr). A non-zero exit is reported in `Result.ExitCode` rather than
as a Go error, so workflows can branch on it.

`tools.go` exposes the `Registry` mapping task name → worker. Adding a tool is
a one-line change.

| Task name | Binary | Skill reference |
|-----------|--------|-----------------|
| `podman`  | podman | [`skills/podman`](../skills/podman/SKILL.md) |
| `kind`    | kind   | [`skills/kind`](../skills/kind/SKILL.md) |
| `k3s`     | k3s    | [`skills/k3s`](../skills/k3s/SKILL.md) |

### Orchestration — `cmd/worker`

Starts a Conductor `TaskRunner`, registers every entry in `tools.Registry`, and
polls until interrupted. Configured via environment variables:

| Variable | Default | Purpose |
|----------|---------|---------|
| `CONDUCTOR_SERVER_URL` | `http://localhost:8080/api` | Conductor API base URL |
| `CONDUCTOR_AUTH_KEY` | _(empty)_ | Orkes application key id (omit for unauthenticated OSS) |
| `CONDUCTOR_AUTH_SECRET` | _(empty)_ | Orkes application key secret |
| `WORKER_BATCH_SIZE` | `1` | tasks polled per cycle, per worker |
| `WORKER_POLL_INTERVAL` | `1s` | poll interval (Go duration) |

```bash
export CONDUCTOR_SERVER_URL=https://play.orkes.io/api
export CONDUCTOR_AUTH_KEY=...    # from Orkes → Access Control → Applications
export CONDUCTOR_AUTH_SECRET=...
go run ./cmd/worker
```

Each tool must also exist as a **task definition** on the server, and workflows
that chain them are registered separately (UI, API, or a follow-up registration
command). Example task input:

```json
{ "args": ["build", "-t", "myapp:dev", "."], "timeoutSeconds": 600 }
```

### Operational excellence — Cortex (Phase 2)

Planned: register each tool-worker in Cortex as a catalog entity and attach
scorecards (ownership, runbook, health). Tracked separately; not yet
implemented.

## Roadmap

1. **Done** — agent-tool workers + Conductor task runner (`internal/tools`, `cmd/worker`).
2. **Done** — Cortex GitOps Tool-Cards for every tool (`cortex.yaml`, `skills/*/cortex.yaml`), stamped with publisher, platform-verifier, version, changelog, and capabilities.
3. **Next** — Cortex scorecards over the catalogued tools (ownership, runbook, health).
4. **Aspiration (soon)** — Agent-Tools operates its own platform and *becomes the
   verifier*, assuming the platform owner's verification duty (signed releases /
   provenance / attestations) instead of delegating it to an external platform.
5. **Deferred** — Temporal, only if a tool needs code-first long-running saga semantics.
