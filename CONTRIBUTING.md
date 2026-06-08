# Contributing to Agent-Tools

Agent-Tools is a **certifying tool-registry**. Contributions add, update, or
deprecate tools and the apps they integrate with. This guide is the community
contribution process; all of it is enforced by CI.

## Principles (non-negotiable)

These come from [`docs/TAXONOMY.md`](docs/TAXONOMY.md) and
[`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md):

1. **No-blocker — everything has an alternative.** Every tool and every platform
   must declare at least one alternative. No single tool or platform may be a
   blocker.
2. **Certification, not publication.** We record `version`/`capabilities` *as
   claimed by the publisher*; verification is the platform owner's duty. We
   never claim to have independently verified a tool.
3. **No fabrication.** Versions must be real published releases (pin via
   `git ls-remote` or the publisher's release page) and the changelog must point
   at exactly that release. Omit unknowns; never invent them.
4. **Editorial is labelled.** `alternatives` and any opinion are Agent-Tools
   editorial, distinct from publisher claims and platform verification.
5. **Optimal, living path.** The roadmap is updated as new information arrives —
   propose path changes when you have better data.

## Add a runnable capability (podman/kind/k3s-style)

1. Add an entry to `internal/tools/tools.go` `Registry` with `Name`, `Binary`,
   real `Version`, `Publisher`, `VerifiedBy`, `Changelog`, `Capabilities`, and
   **`Alternatives` (≥1)**, plus a `Worker`.
2. Add a worker func (`func Foo(t *model.Task) (interface{}, error) { return run("foo", t) }`).
3. Add a skill at `skills/<tool>/SKILL.md`.
4. Add a Tool-Card at `skills/<tool>/cortex.yaml` from
   [`templates/tool-card.cortex.yaml`](templates/tool-card.cortex.yaml).

## Catalog an integration surface (SDK / CLI / MCP / API / GitOps)

Add it under the right app in [`catalog/apps.yaml`](catalog/apps.yaml) with
`type` (from the taxonomy), `provenance` (`official`/`community`), and `source`.
SDK entries need a `language`. Each app must keep ≥1 `alternatives`.

## Pin a real version

```bash
git ls-remote --tags --refs https://github.com/<org>/<repo>.git \
  | awk -F/ '{print $NF}' | grep -E '^v?[0-9]+\.' | sort -V | tail -3
```

Use the exact tag and set `changelog` to that release's page.

## Validate locally before opening a PR

```bash
go build ./... && go vet ./... && go test ./...   # Go invariants
python scripts/validate_registry.py               # taxonomy + no-blocker rules
```

CI runs the same checks and **fails on any rule violation** — do not break rules.

## PR checklist

- [ ] Real version, changelog points at exactly that release
- [ ] `alternatives` declared (≥1) — no-blocker rule
- [ ] `type`/`provenance` use the controlled vocabulary
- [ ] Editorial fields labelled as such
- [ ] `go test ./...` and `scripts/validate_registry.py` pass locally
