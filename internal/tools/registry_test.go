package tools

import "testing"

// TestRegistryInvariants enforces the registry rules from docs/TAXONOMY.md in
// code, so `go test` fails the build if an entry breaks them. The most load-
// bearing is the no-blocker rule: every tool must declare an alternative.
func TestRegistryInvariants(t *testing.T) {
	if len(Registry) == 0 {
		t.Fatal("Registry is empty")
	}
	for name, tool := range Registry {
		if tool.Name != name {
			t.Errorf("%s: Name %q does not match registry key", name, tool.Name)
		}
		if tool.Binary == "" {
			t.Errorf("%s: missing Binary", name)
		}
		if tool.Version == "" {
			t.Errorf("%s: missing Version", name)
		}
		if tool.Publisher == "" {
			t.Errorf("%s: missing Publisher", name)
		}
		if tool.VerifiedBy == "" {
			t.Errorf("%s: missing VerifiedBy (platform owner)", name)
		}
		if tool.Changelog == "" {
			t.Errorf("%s: missing Changelog (certification evidence)", name)
		}
		if len(tool.Capabilities) == 0 {
			t.Errorf("%s: must declare at least one capability", name)
		}
		// Rule 1 — no-blocker: everything has an alternative.
		if len(tool.Alternatives) == 0 {
			t.Errorf("%s: no-blocker rule violated — must declare at least one alternative", name)
		}
		if tool.Worker == nil {
			t.Errorf("%s: missing Worker", name)
		}
	}
}
