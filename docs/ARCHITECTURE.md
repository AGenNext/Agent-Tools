# Agent-Tools Architecture

This repo turns container / Kubernetes capabilities (podman, kind, k3s) into
**orchestrated agent tools**. The system is layered so each concern is owned by
the right component, and every SDK in the stack stays available — none is
removed.

```
┌──────────────────────────────────────────────────────────┐
│  OPERATIONAL EXCELLENCE   →  Cortex (IDP)                  │
│  catalog every tool as an entity, ownership,              │
│  scorecards (runbook? healthy? on-call?), standards       │
├──────────────────────────────────────────────────────────┤
│  ORCHESTRATION / EXECUTION                                 │
│    • Orkes Conductor  — declarative workflows, polled      │
│                         workers, visual run history        │
│    • Temporal         — code-first durable execution,      │
│                         long-running stateful sagas        │
├──────────────────────────────────────────────────────────┤
│  CAPABILITIES (agent tools)  →  Go workers / activities    │
│  podman / kind / k3s — run the CLI, return stdout/exit     │
└──────────────────────────────────────────────────────────┘
```

## SDKs

All SDKs below are part of the stack and remain supported. The same shell-free
tool runner (`internal/tools`) backs every binding, so a tool is written once
and exposed through whichever engine a workflow needs.

| SDK / platform | Module / source | Role | Status |
|---|---|---|---|
| **Orkes Conductor Go SDK** | `github.com/conductor-sdk/conductor-go` | Orchestration + execution: each tool is a registered Conductor worker (`cmd/worker`). | Implemented |
| **Temporal Go SDK** | `go.temporal.io/sdk` | Code-first durable execution: each tool is exposed as a Temporal Activity for long-running, stateful, retried/heartbeated sagas. | Planned (same runner, Activity binding) |
| **Cortex** | GitOps (`cortex.yaml`) + REST API | Operational excellence / IDP: catalogs the tools and runs scorecards. Never executes work. | Implemented (GitOps Tool-Cards) |
| **Orkes Conductor metadata API** | via the Conductor Go SDK | Registers task & workflow definitions on the server. | Planned |

### Why both Conductor and Temporal (not either/or)

They cover different sweet spots, so both are kept:

- **Conductor** — declarative, polyglot, polled workers, visual run history;
  ideal for composing tools into visible pipelines and adding tools cheaply.
- **Temporal** — code-first workflows-as-Go with strong durability primitives
  (timers, signals, heartbeats); ideal when a single tool operation is a
  long-running, stateful saga (hours–days, human-in-the-loop).

A Conductor workflow can delegate a heavy step to a Temporal workflow when that
step needs Temporal's durability — the two compose rather than compete.

**Cortex** is not a workflow engine at all — it never executes work. It is the
Internal Developer Portal layer that catalogs the tools and runs scorecards for
operational excellence.

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

**Integration mode decides who validates.** *How* we consume a tool changes
where validation responsibility sits:

| Integration mode | Example | Who validates |
|---|---|---|
| **API / MCP** (called over the wire) | Conductor/Temporal/Cortex MCP servers, REST APIs | External official surface — verification stays with the **platform owner**. We ship none of their code. |
| **SDK** (compiled into our build) | `conductor-sdk/conductor-go` in `go.mod` | Enters **our supply chain** — validation responsibility shifts to **Agent-Tools** (pin in `go.sum`, Dependabot, vuln scanning). |

This is why the repo runs Dependabot and why we bump vulnerable transitive deps
(e.g. `protobuf` → v1.34.2): we *use* the Conductor Go SDK, so we own its
validation. MCP and API entries are official external tools and carry no such
build-time obligation.

**Design preference: MCP / API over SDK.** Wherever a capability is reachable
over the wire, prefer the platform's MCP server or API to embedding its SDK.
This (a) keeps validation with the platform owner, (b) shrinks our supply chain,
and (c) keeps the audit trail clean: **the action is logged with Agent-Tools as
the actor and the platform (GitHub, Conductor, …) as the tool**, rather than us
re-shipping their code. Reach for an SDK only when no MCP/API path exists or the
work must run in-process.

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
3. **Done** — complete app integration catalog: every SDK (all languages), CLI,
   and MCP server for Conductor, Temporal, and Cortex catalogued in
   [`catalog/apps.yaml`](../catalog/apps.yaml) / [`docs/CATALOG.md`](CATALOG.md).
4. **Next** — Temporal Go SDK binding: expose the same tool runner as Temporal
   Activities (kept as a first-class engine, not deferred).
5. **Next** — Cortex scorecards over the catalogued tools (ownership, runbook, health).
6. **Aspiration (soon)** — Agent-Tools operates its own platform and *becomes the
   verifier*, assuming the platform owner's verification duty (signed releases /
   provenance / attestations) instead of delegating it to an external platform.
