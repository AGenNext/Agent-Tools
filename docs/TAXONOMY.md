# Tool Taxonomy & Vocabulary

This is the schema authority for the Agent-Tools registry. Every entry — whether
a runnable capability (`internal/tools` + `skills/*/cortex.yaml`) or a catalogued
integration surface (`catalog/apps.yaml`) — uses this controlled vocabulary. The
[validator](../scripts/validate_registry.py) enforces it in CI.

## Layers

| Layer | Meaning | Where |
|---|---|---|
| `capability` | A tool that does work (a CLI-backed worker/activity). | `internal/tools`, `skills/*` |
| `orchestration` | Engines that schedule/execute workflows over capabilities. | Conductor, Temporal |
| `operational-excellence` | Catalog / scorecards / governance. | Cortex (IDP) |

## Tool kinds — `type`

Controlled set (`scripts/validate_registry.py` rejects anything else):

| `type` | Meaning |
|---|---|
| `capability` | A binary-backed tool the orchestrator executes (podman, kind, k3s). |
| `sdk` | A client library for a platform, qualified by `language`. |
| `cli` | A command-line client for a platform. |
| `mcp` | A Model Context Protocol server. |
| `gitops` | A Git-synced descriptor mechanism (e.g. `cortex.yaml`). |
| `api` | A network API surface (e.g. REST). |

## Domains — `category`

`container-runtime` · `kubernetes` · `orchestration` · `durable-execution` · `idp`

## Trust & provenance vocabulary

These mirror the [trust model](ARCHITECTURE.md#trust--certification-model):

| Field | Meaning | Source of truth |
|---|---|---|
| `publisher` | Who ships the tool and makes the claims. | the publisher |
| `verified-by` | The **platform owner** whose duty it is to verify the claims (e.g. GitHub — signed tags / releases / provenance). | platform owner |
| `version` | Pinned release, **as claimed by the publisher**. | publisher changelog |
| `changelog` | Platform-hosted evidence URL for exactly `version`. | publisher/platform |
| `capabilities` | Feature set, **as claimed by the publisher**. | publisher |
| `provenance` | `official` (vendor-maintained) or `community` (third-party). | Agent-Tools assessment |
| `certification` | Disclaimer: Agent-Tools references claims, does not warrant/verify. | Agent-Tools |
| `alternatives` | Substitute tools/platforms — **Agent-Tools editorial**, not a publisher claim. | Agent-Tools |
| `status` | `integrated` · `catalogued` · `planned` · `deprecated`. | Agent-Tools |

## Rules (invariants — enforced, do not break)

1. **No-blocker / everything has an alternative.** Every capability Tool-Card and
   every catalogued platform (`app`) MUST declare at least one `alternatives`
   entry. No single tool or platform may be a blocker.
2. **Certification, not publication.** `version` and `capabilities` are recorded
   as *claimed by the publisher*; verification is the platform owner's duty.
   Agent-Tools never asserts it independently verified a tool.
3. **No fabrication.** `version` must be a real published release (pin via
   `git ls-remote` / the publisher's release page) and `changelog` must point at
   exactly that release. Unknown values are omitted, never invented.
4. **Controlled vocabulary.** `type` must be one of the kinds above; `provenance`
   must begin with `official` or `community`.
5. **Editorial fields are labelled.** `alternatives` (and any opinion) are marked
   as Agent-Tools editorial, distinct from publisher claims and platform
   verification.

Violating any rule fails CI.
