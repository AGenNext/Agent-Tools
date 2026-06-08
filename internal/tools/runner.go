package tools

import (
	"context"
	"fmt"

	"github.com/conductor-sdk/conductor-go/sdk/model"
)

// run adapts a Conductor task to the engine-agnostic executor (see exec.go). It
// decodes the recognized task inputs and delegates to Run.
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

	in := Invocation{Args: args}
	if v, ok := t.InputData["timeoutSeconds"]; ok {
		if secs, ok := toInt(v); ok {
			in.TimeoutSeconds = secs
		}
	}
	if stdin, ok := t.InputData["stdin"].(string); ok {
		in.Stdin = stdin
	}

	return Run(context.Background(), binary, in)
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
