package tools

import "github.com/conductor-sdk/conductor-go/sdk/model"

// Tool describes a single agent tool in the registry: how Conductor addresses
// it (Name -> task name), which binary it drives, the Version and Capabilities
// CERTIFIED for it, and the Worker that executes it. This struct is the in-code
// source of truth that the Cortex Tool-Cards (skills/<tool>/cortex.yaml) mirror.
//
// Three-way trust split:
//   - the Publisher MAKES the claims (Version, Capabilities) in its Changelog;
//   - the hosting platform OWNER (VerifiedBy, e.g. GitHub) VERIFIES those claims
//     via its own primitives — signed tags, verified releases, build provenance
//     / attestations — because authenticity is, everywhere, the platform
//     owner's responsibility and duty;
//   - Agent-Tools TRUSTS the platform owner as the verification authority and
//     CERTIFIES only that this entry faithfully references the platform-hosted
//     source. It does not publish, warrant, or independently verify the tool.
//
// Future: when Agent-Tools operates its own platform, it becomes the platform
// owner and assumes the verification duty itself, rather than delegating it.
type Tool struct {
	Name    string
	Binary  string
	Version string
	// Publisher identifies who ships the tool and whose claims this entry
	// references. It is free-form (a GitHub org, a vendor name, a GitLab group,
	// ...) so the registry never assumes a single source.
	Publisher string
	// VerifiedBy names the hosting platform responsible for verifying the
	// publisher's claims (e.g. "GitHub" — signed tags / verified releases /
	// build provenance). Verification is the platform's job, not the registry's.
	VerifiedBy string
	// Changelog is the publisher's release-notes URL for the pinned Version,
	// hosted on VerifiedBy's platform. It is the source the Version +
	// Capabilities claims are certified against; it must point at the release
	// for exactly Version, so a version bump and its changelog move together.
	Changelog string
	// Capabilities are AS CLAIMED BY THE PUBLISHER for this Version, not
	// independently verified by Agent-Tools.
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
		Version:      "v5.8.2",
		Publisher:    "containers (Red Hat)",
		VerifiedBy:   "GitHub",
		Changelog:    "https://github.com/containers/podman/releases/tag/v5.8.2",
		Capabilities: []string{"build", "run", "images", "pods", "kube-generate", "push"},
		Worker:       Podman,
	},
	"kind": {
		Name:         "kind",
		Binary:       "kind",
		Version:      "v0.32.0",
		Publisher:    "kubernetes-sigs",
		VerifiedBy:   "GitHub",
		Changelog:    "https://github.com/kubernetes-sigs/kind/releases/tag/v0.32.0",
		Capabilities: []string{"cluster-create", "cluster-delete", "load-image", "multi-node", "export-kubeconfig"},
		Worker:       Kind,
	},
	"k3s": {
		Name:         "k3s",
		Binary:       "k3s",
		Version:      "v1.36.1+k3s1",
		Publisher:    "k3s-io (SUSE/Rancher)",
		VerifiedBy:   "GitHub",
		Changelog:    "https://github.com/k3s-io/k3s/releases/tag/v1.36.1+k3s1",
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
