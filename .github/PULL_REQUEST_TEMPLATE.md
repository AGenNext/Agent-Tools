<!-- Agent-Tools is a certifying tool-registry. Respect every rule (docs/TAXONOMY.md). -->

## What does this change?

<!-- New tool? Version bump? Catalog/SDK/CLI/MCP entry? Docs? -->

## Rule checklist (CI enforces these)

- [ ] **No-blocker:** every added/changed tool or platform declares ≥1 `alternatives`
- [ ] **No fabrication:** `version` is a real published release; `changelog` points at exactly that release
- [ ] **Certification, not publication:** version/capabilities recorded as *publisher-claimed*; not asserted as independently verified
- [ ] **Vocabulary:** `type`/`provenance` use the controlled sets in `docs/TAXONOMY.md`
- [ ] **Editorial labelled:** `alternatives`/opinions marked as Agent-Tools editorial
- [ ] Ran locally: `go test ./...` and `python scripts/validate_registry.py` pass
