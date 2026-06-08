package tools

import (
	"context"
	"strings"
	"testing"
)

// TestRunTimeoutSurfacesError guards the fix for the masked-timeout bug: when a
// process is killed by the per-invocation timeout it returns an *exec.ExitError,
// and Run must report that as a timeout error rather than a successful Result.
func TestRunTimeoutSurfacesError(t *testing.T) {
	res, err := Run(context.Background(), "sleep", Invocation{
		Args:           []string{"10"},
		TimeoutSeconds: 1,
	})
	if err == nil {
		t.Fatalf("expected timeout error, got success: %+v", res)
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("expected a timeout error, got: %v", err)
	}
}

// TestRunNonZeroExit confirms an ordinary non-zero exit is reported via the
// Result (exit code), not as a Go error.
func TestRunNonZeroExit(t *testing.T) {
	res, err := Run(context.Background(), "false", Invocation{})
	if err != nil {
		t.Fatalf("unexpected error for non-zero exit: %v", err)
	}
	if res.ExitCode == 0 {
		t.Errorf("expected non-zero ExitCode, got 0")
	}
}
