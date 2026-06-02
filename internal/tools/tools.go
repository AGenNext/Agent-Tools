package tools

import "github.com/conductor-sdk/conductor-go/sdk/model"

// Tool describes a single agent tool in the registry: how Conductor addresses
// it (Name -> task name), which binary it drives, the contract Version of this
// Tool-Card, the Capabilities it advertises, and the Worker that executes it.
// This struct is the in-code source of truth that the Cortex Tool-Cards
// (skills/<tool>/cortex.yaml) mirror.
type Tool struct {
	Name         string
	Binary       string
	Version      string
	Capabilities []string
	Worker       model.ExecuteTaskFunction
}

// Registry maps a Conductor task name to its Tool. The cmd/worker entrypoint
// iterates this map and registers each entry as an agent tool, so adding a tool
// is a single entry here (plus its Tool-Card under skills/<tool>/cortex.yaml).
var Registry = map[string]Tool{
	"podman": {
		Name:         "podman",
		Binary:       "podman",
		Version:      "1.0.0",
		Capabilities: []string{"build", "run", "images", "pods", "kube-generate", "push"},
		Worker:       Podman,
	},
	"kind": {
		Name:         "kind",
		Binary:       "kind",
		Version:      "1.0.0",
		Capabilities: []string{"cluster-create", "cluster-delete", "load-image", "multi-node", "export-kubeconfig"},
		Worker:       Kind,
	},
	"k3s": {
		Name:         "k3s",
		Binary:       "k3s",
		Version:      "1.0.0",
		Capabilities: []string{"server", "agent", "kubectl", "cluster-init-ha", "crictl"},
		Worker:       K3s,
	},
}

// Podman executes the `podman` CLI. See run for the recognized task inputs.
func Podman(t *model.Task) (interface{}, error) { return run("podman", t) }

// Kind executes the `kind` (Kubernetes IN Docker) CLI.
func Kind(t *model.Task) (interface{}, error) { return run("kind", t) }

// K3s executes the `k3s` CLI.
func K3s(t *model.Task) (interface{}, error) { return run("k3s", t) }
