#!/usr/bin/env python3
"""Validate the Agent-Tools registry against the taxonomy in docs/TAXONOMY.md.

Enforces the invariants (see docs/TAXONOMY.md "Rules"):
  - no-blocker: every Tool-Card and every catalogued app declares >= 1 alternative
  - controlled vocabulary for `type` / `provenance`
  - certification fields present (publisher, verified-by, version, changelog,
    capabilities, certification) on Tool-Cards
  - catalogued tools carry type/provenance/source

Exits non-zero (failing CI) on any violation.
"""
from __future__ import annotations

import glob
import sys

try:
    import yaml
except ImportError:
    sys.exit("PyYAML is required: pip install pyyaml")

TOOL_KINDS = {"capability", "sdk", "cli", "mcp", "gitops", "api"}
PROVENANCE_PREFIXES = ("official", "community")

errors: list[str] = []


def err(where: str, msg: str) -> None:
    errors.append(f"{where}: {msg}")


def load(path: str):
    with open(path) as fh:
        return yaml.safe_load(fh)


def csv_nonempty(value) -> bool:
    """A comma-separated scalar (or list) with at least one real item."""
    if value is None:
        return False
    if isinstance(value, list):
        return len([x for x in value if str(x).strip()]) > 0
    return len([p for p in str(value).split(",") if p.strip()]) > 0


def validate_catalog(path: str = "catalog/apps.yaml") -> None:
    data = load(path)
    apps = (data or {}).get("apps")
    if not apps:
        err(path, "no `apps` list found")
        return
    for app in apps:
        name = app.get("name", "<unnamed>")
        where = f"{path} [{name}]"
        for field in ("name", "role", "publisher", "verified-by"):
            if not app.get(field):
                err(where, f"missing required app field `{field}`")
        # Rule 1: no-blocker for platforms.
        if not csv_nonempty(app.get("alternatives")):
            err(where, "no-blocker rule: app must declare >= 1 `alternatives`")
        tools = app.get("tools") or []
        if not tools:
            err(where, "app has no `tools`")
        for tool in tools:
            tname = tool.get("name", "<unnamed>")
            twhere = f"{where} -> {tname}"
            ttype = tool.get("type")
            if ttype not in TOOL_KINDS:
                err(twhere, f"invalid `type` {ttype!r}; allowed: {sorted(TOOL_KINDS)}")
            prov = str(tool.get("provenance", ""))
            if not prov.startswith(PROVENANCE_PREFIXES):
                err(twhere, f"`provenance` must start with {PROVENANCE_PREFIXES}; got {prov!r}")
            if not tool.get("source"):
                err(twhere, "missing `source` URL")
            if ttype == "sdk" and not tool.get("language"):
                err(twhere, "sdk entry must specify `language`")


REQUIRED_CARD_META = (
    "publisher",
    "verified-by",
    "version",
    "changelog",
    "capabilities",
    "certification",
    "alternatives",
)


def validate_tool_cards(pattern: str = "skills/*/cortex.yaml") -> None:
    paths = sorted(glob.glob(pattern))
    if not paths:
        err(pattern, "no Tool-Cards found")
    for path in paths:
        info = (load(path) or {}).get("info") or {}
        if not info.get("x-cortex-tag"):
            err(path, "missing info.x-cortex-tag")
        meta = info.get("x-cortex-custom-metadata") or {}
        for field in REQUIRED_CARD_META:
            if field not in meta:
                err(path, f"Tool-Card missing custom-metadata `{field}`")
        # Rule 1: no-blocker for capabilities.
        if not csv_nonempty(meta.get("alternatives")):
            err(path, "no-blocker rule: Tool-Card must declare >= 1 `alternatives`")
        if not csv_nonempty(meta.get("capabilities")):
            err(path, "Tool-Card must declare >= 1 `capabilities`")
        # Rule 3 (light): a version implies a changelog pointing somewhere.
        if meta.get("version") and not str(meta.get("changelog", "")).startswith("http"):
            err(path, "`version` present but `changelog` is not a URL")


def main() -> int:
    validate_catalog()
    validate_tool_cards()
    if errors:
        print("Registry validation FAILED:\n")
        for e in errors:
            print(f"  - {e}")
        print(f"\n{len(errors)} violation(s).")
        return 1
    print("Registry validation passed (taxonomy + no-blocker rules).")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
