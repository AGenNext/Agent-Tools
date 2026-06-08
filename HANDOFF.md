# Session Handoff

A durable brief so a future session (or person) can resume without re-deriving
anything. Committed because session memory does not carry over; only what's in a
repo does.

## State of Agent-Tools

- **PR #2** (`claude/zealous-cray-3RoX4` → `main`) is **ready for review** and
  contains the whole registry: capabilities (podman/kind/k3s workers), the
  catalog (`catalog/apps.yaml`), trust model + taxonomy + rules, governance,
  and the consolidated security posture.
- Orientation lives in [`AGENTS.md`](AGENTS.md); live status/next loops in
  [`docs/PROJECT.md`](docs/PROJECT.md). Read those first.
- Latest fix: `internal/tools/exec.go` now surfaces context timeouts before
  treating them as ordinary exits (Codex review on PR #2, resolved).

## Active task: build the skill registry in `AGenNext/Agent-Skills`

Agent-Skills is a **separate repo** and was **out of this session's access**
(only `agennext/agent-tools` granted; an agent cannot grant itself access — a
human must add the repo to the session, or start a session scoped to it).

What was learned by reading it publicly:

- It is **schema-driven YAML**, not Markdown. Authoring skills live at
  `skills/<category>/<id>.yaml` (e.g. `skills/cloud/k3s.install.yaml`).
- DB-record contract: `schemas/skill.schema.json` (`skill_source_v1`). Required:
  `type`, `skill_id`, `name`, `version`, `lifecycle_status`, `category`, `risk`,
  `runtime_profiles`, `required_tools`, `approval_required`, `inputs`,
  `outputs`, and the graph-ownership consts (`graph_context_owner: Agent-Graph`,
  `semantic_class: CreativeWork`, `graph_ownership: Agent-Graph`,
  `does_not_define_graph_schema: true`).
- Repo boundary: Agent-Skills owns skill definitions/metadata/contracts only;
  composition/runtime/SDKs live in Agent-Team / Agent-Runtime / AgentKube.
- Planned cloud skills: `cloud.inspect_server`, `linux.harden_node`,
  `k3s.install`, `kubernetes.deploy_manifest`, `surrealdb.deploy`,
  `security.preflight`, `compliance.collect_evidence`, `lifecycle.update_state`.

### Validated draft to carry over (not yet pushable — no access)

This `k3s.install` record was validated green against the real
`schemas/skill.schema.json`:

```yaml
type: ag:Skill
skill_id: k3s.install
name: Install k3s
version: 0.1.0
schema_version: skill_source_v1
lifecycle_status: draft
category: cloud
risk: high
description: Install k3s on a target host over SSH; return the kubeconfig and node status.
runtime_profiles: [k8smicro]
required_tools: [ssh.run_script, kubernetes.get_nodes]
required_secrets: [ssh_private_key]
approval_required: true
inputs:
  target_host: string
  ssh_user: string
outputs:
  kubeconfig: secret_ref
  node_status: object
graph_context_owner: Agent-Graph
semantic_class: CreativeWork
graph_ownership: Agent-Graph
does_not_define_graph_schema: true
source_repository: AGenNext/Agent-Skills
source_path: skills/cloud/k3s.install.yaml
source_format: yaml
```

### Next loop (once Agent-Skills access is granted)

Apply the same governance pattern we built here:
1. Make each `skills/<category>/*.yaml` conform to (or map cleanly to)
   `schemas/skill.schema.json`.
2. Add a `validate` step (jsonschema) + CI mirroring `Agent-Tools` (GitHub-native
   security, one-task-one-tool).
3. Generate a skill registry index from the skill files (single source of truth,
   no drift) — analogous to this repo's card↔registry test.

## The only human-only step

Granting repo access. Everything else is automatable from inside a properly
scoped session.
