# App Integration Catalog

In Agent-Tools **everything is a tool** — every way to integrate with an app
(SDK, CLI, MCP server, GitOps, REST API) is catalogued as a tool. The
machine-readable source of truth is [`catalog/apps.yaml`](../catalog/apps.yaml);
this page summarizes it.

Each entry records what the vendor **publishes**, its **provenance** (official
vs community), and the **platform owner** responsible for verifying the source
(per the [trust model](ARCHITECTURE.md#trust--certification-model)). Agent-Tools
certifies the entry references the vendor source — it does not verify or warrant
the tool.

## Orkes Conductor — orchestration + execution
Publisher: conductor-oss / Orkes · Verified by: GitHub

| Tool | Type | Language | Provenance |
|---|---|---|---|
| Go SDK *(integrated)* | SDK | Go | official |
| Python SDK | SDK | Python | official |
| Java SDK | SDK | Java | official |
| JavaScript / TypeScript SDK | SDK | JS/TS | official |
| C# SDK | SDK | C# | official |
| Clojure SDK | SDK | Clojure | official |
| Conductor MCP Server | MCP | — | official |

## Temporal — code-first durable execution
Publisher: Temporal Technologies · Verified by: GitHub

| Tool | Type | Language | Provenance |
|---|---|---|---|
| Go SDK | SDK | Go | official |
| Java SDK | SDK | Java | official |
| TypeScript SDK | SDK | TypeScript | official |
| Python SDK | SDK | Python | official |
| .NET SDK | SDK | .NET | official |
| PHP SDK | SDK | PHP | official |
| Ruby SDK | SDK | Ruby | official |
| Temporal CLI | CLI | — | official |
| Temporal MCP Server | MCP | — | community (Temporal Code Exchange) |

## Cortex (cortex.io IDP) — operational excellence
Publisher: Cortex (cortexapps) · Verified by: GitHub

| Tool | Type | Provenance |
|---|---|---|
| GitOps descriptors (`cortex.yaml`) *(integrated)* | GitOps | official |
| REST API | API | official |
| Cortex MCP Server | MCP | official |

> **Name disambiguation:** this Cortex is the **cortex.io Internal Developer
> Portal** (cortexapps). *Snowflake Cortex* and *Palo Alto Cortex* are unrelated
> products that share the name and are excluded. Cortex.io publishes no language
> SDKs (integrate via REST API + GitOps); no official cortex.io CLI is confirmed.
