package tools

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The Go Registry (tools.go) is the single source of truth. The Cortex
// Tool-Cards (skills/<tool>/cortex.yaml) mirror it for GitOps. This test
// enforces that mirror so the two can never drift — one task, one owner.

// cardMeta extracts a custom-metadata value from a Tool-Card, stripping any
// inline comment. Kept deliberately dependency-free (no YAML library) to avoid
// adding to the supply chain for a test.
func cardMeta(t *testing.T, tool, key string) string {
	t.Helper()
	path := filepath.Join("..", "..", "skills", tool, "cortex.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read card %s: %v", path, err)
	}
	re := regexp.MustCompile(`(?m)^    ` + regexp.QuoteMeta(key) + `:\s*(.*)$`)
	m := re.FindStringSubmatch(string(data))
	if m == nil {
		t.Fatalf("%s: custom-metadata key %q not found in card", tool, key)
	}
	v := m[1]
	if i := strings.Index(v, " #"); i >= 0 {
		v = v[:i] // strip inline comment
	}
	return strings.TrimSpace(v)
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func equalSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestToolCardsMatchRegistry(t *testing.T) {
	for name, tool := range Registry {
		if got := cardMeta(t, name, "version"); got != tool.Version {
			t.Errorf("%s: card version %q != registry %q", name, got, tool.Version)
		}
		if got := cardMeta(t, name, "publisher"); got != tool.Publisher {
			t.Errorf("%s: card publisher %q != registry %q", name, got, tool.Publisher)
		}
		if got := splitCSV(cardMeta(t, name, "capabilities")); !equalSlice(got, tool.Capabilities) {
			t.Errorf("%s: card capabilities %v != registry %v", name, got, tool.Capabilities)
		}
		if got := splitCSV(cardMeta(t, name, "alternatives")); !equalSlice(got, tool.Alternatives) {
			t.Errorf("%s: card alternatives %v != registry %v", name, got, tool.Alternatives)
		}
	}
}
