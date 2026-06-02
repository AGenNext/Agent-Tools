package tools

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Result is the standard output shape returned by every agent-tool execution,
// regardless of engine (Conductor worker or Temporal activity). It is
// serialized into the engine's task/activity output so downstream steps can
// branch on ExitCode or parse Stdout.
type Result struct {
	Binary   string   `json:"binary"`
	Args     []string `json:"args"`
	ExitCode int      `json:"exitCode"`
	Stdout   string   `json:"stdout"`
	Stderr   string   `json:"stderr"`
}

// Invocation is the engine-agnostic input for running a tool. Conductor decodes
// it from task InputData; Temporal passes it as the activity argument.
type Invocation struct {
	Args           []string `json:"args"`
	TimeoutSeconds int      `json:"timeoutSeconds"`
	Stdin          string   `json:"stdin"`
}

// DefaultTimeout bounds any single tool invocation that doesn't specify its own
// TimeoutSeconds.
const DefaultTimeout = 5 * time.Minute

// Run executes a fixed binary with the supplied args. We invoke the binary
// directly (no shell) so caller-supplied args cannot inject additional
// commands. A non-zero exit is NOT treated as a Go error — the exit code is
// surfaced in the Result so callers can decide how to react; a Go error is
// reserved for "couldn't run at all" cases (binary missing, timeout). The
// passed ctx is honored (e.g. activity cancellation) in addition to the
// per-invocation timeout.
func Run(ctx context.Context, binary string, in Invocation) (Result, error) {
	timeout := DefaultTimeout
	if in.TimeoutSeconds > 0 {
		timeout = time.Duration(in.TimeoutSeconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, binary, in.Args...)
	if in.Stdin != "" {
		cmd.Stdin = strings.NewReader(in.Stdin)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitCode := 0
	if runErr := cmd.Run(); runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			if ctx.Err() == context.DeadlineExceeded {
				return Result{}, fmt.Errorf("%s timed out after %s", binary, timeout)
			}
			return Result{}, fmt.Errorf("failed to run %s: %w (stderr: %s)", binary, runErr, strings.TrimSpace(stderr.String()))
		}
	}

	return Result{
		Binary:   binary,
		Args:     in.Args,
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}
