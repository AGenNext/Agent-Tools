# Agent-Tools

A collection of [Claude Code Agent Skills](https://docs.claude.com/en/docs/claude-code/skills)
for working with container and local-Kubernetes tooling.

## Skills

| Skill | Use it for | Don't use it for |
| --- | --- | --- |
| [`k3s`](skills/k3s/SKILL.md) | A real, persistent single-/multi-node Kubernetes cluster on a host, VM, edge device, or CI runner (single binary, no Docker). | Throwaway clusters inside containers — use `kind`. |
| [`kind`](skills/kind/SKILL.md) | Disposable Kubernetes clusters running *inside* Docker/Podman containers for local dev, testing controllers, and CI. | Long-lived / edge / production nodes — use `k3s`. |
| [`podman`](skills/podman/SKILL.md) | Building and running OCI containers, images, and pods with a daemonless, rootless-capable, Docker-compatible engine. | Orchestrating a multi-node cluster — use `kind`/`k3s`. |

These tools compose: **podman** (or Docker) can back a **kind** cluster, and
`podman kube generate` produces manifests you can apply to a **kind** or
**k3s** cluster.

## Orchestrator

These tools are also exposed as **Orkes Conductor** workers so an agent can chain
them into durable, multi-step workflows. `cmd/worker` registers each tool
(`podman`, `kind`, `k3s`) as a Conductor task worker. See
[`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) for the full layering
(Conductor for orchestration, Go workers as capabilities, Cortex for
operational excellence — and why Temporal is deferred).

```bash
export CONDUCTOR_SERVER_URL=https://play.orkes.io/api
export CONDUCTOR_AUTH_KEY=...
export CONDUCTOR_AUTH_SECRET=...
go run ./cmd/worker
```

## Layout

```
skills/
  k3s/SKILL.md
  kind/SKILL.md
  podman/SKILL.md
```

Each `SKILL.md` has YAML frontmatter (`name`, `description`) describing when the
skill applies, followed by practical, agent-oriented usage instructions.
