package tools

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/conductor-sdk/conductor-go/sdk/model"
)

// Result is the standard output shape returned by every agent-tool worker.
// It is serialized into the Conductor task output, so downstream tasks in a
// workflow can branch on ExitCode or parse Stdout.
type Result struct {
	Binary   string   `json:"binary"`
	Args     []string `json:"args"`
	ExitCode int      `json:"exitCode"`
	Stdout   string   `json:"stdout"`
	Stderr   string   `json:"stderr"`
}

// DefaultTimeout bounds any single tool invocation that doesn't specify its
// own "timeoutSeconds" input.
const DefaultTimeout = 5 * time.Minute

// run executes a fixed binary with the args supplied in the task input. We
// invoke the binary directly (no shell) so user-supplied args cannot inject
// additional commands. A non-zero exit is NOT treated as a Go error — the exit
// code is surfaced in the Result so workflows can decide how to react; a Go
// error is reserved for "couldn't run at all" cases (binary missing, timeout,
// bad input).
//
// Recognized task inputs:
//
//	args           []string  arguments passed to the binary (required)
//	timeoutSeconds number    per-invocation timeout override (optional)
//	stdin          string    data piped to the process stdin (optional)
func run(binary string, t *model.Task) (interface{}, error) {
	args, err := stringSlice(t.InputData["args"])
	if err != nil {
		return nil, fmt.Errorf("invalid %q input: %w", "args", err)
	}

	timeout := DefaultTimeout
	if v, ok := t.InputData["timeoutSeconds"]; ok {
		if secs, ok := toInt(v); ok && secs > 0 {
			timeout = time.Duration(secs) * time.Second
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, binary, args...)
	if stdin, ok := t.InputData["stdin"].(string); ok && stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitCode := 0
	if runErr := cmd.Run(); runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			// The process ran but returned non-zero. Report it via the Result.
			exitCode = exitErr.ExitCode()
		} else {
			// Couldn't start (binary not found), timed out, etc.
			if ctx.Err() == context.DeadlineExceeded {
				return nil, fmt.Errorf("%s timed out after %s", binary, timeout)
			}
			return nil, fmt.Errorf("failed to run %s: %w (stderr: %s)", binary, runErr, strings.TrimSpace(stderr.String()))
		}
	}

	return Result{
		Binary:   binary,
		Args:     args,
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}

// stringSlice coerces a JSON-decoded value into []string. Conductor decodes
// arrays as []interface{}, so we handle that case as well as a native
// []string and a single string.
func stringSlice(v interface{}) ([]string, error) {
	switch val := v.(type) {
	case nil:
		return nil, nil
	case []string:
		return val, nil
	case string:
		if val == "" {
			return nil, nil
		}
		return []string{val}, nil
	case []interface{}:
		out := make([]string, 0, len(val))
		for i, e := range val {
			s, ok := e.(string)
			if !ok {
				return nil, fmt.Errorf("element %d is %T, want string", i, e)
			}
			out = append(out, s)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("got %T, want array of strings", v)
	}
}

// toInt coerces a JSON-decoded number (float64) or string into an int.
func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	case string:
		var i int
		if _, err := fmt.Sscanf(n, "%d", &i); err == nil {
			return i, true
		}
	}
	return 0, false
}
