package tools

import "github.com/conductor-sdk/conductor-go/sdk/model"

// Registry maps a Conductor task (reference) name to the worker that executes
// it. The cmd/worker entrypoint iterates this map and registers each entry as
// an agent tool, so adding a new tool is a one-line change here.
var Registry = map[string]model.ExecuteTaskFunction{
	"podman": Podman,
	"kind":   Kind,
	"k3s":    K3s,
}

// Podman executes the `podman` CLI. See run for the recognized task inputs.
func Podman(t *model.Task) (interface{}, error) { return run("podman", t) }

// Kind executes the `kind` (Kubernetes IN Docker) CLI.
func Kind(t *model.Task) (interface{}, error) { return run("kind", t) }

// K3s executes the `k3s` CLI.
func K3s(t *model.Task) (interface{}, error) { return run("k3s", t) }
