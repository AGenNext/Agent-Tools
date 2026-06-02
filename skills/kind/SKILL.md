---
name: kind
description: >-
  Create and manage disposable Kubernetes clusters that run inside Docker or
  Podman containers using kind (Kubernetes IN Docker). Use for local
  development, testing controllers/operators, loading locally-built images, and
  CI where a throwaway multi-node cluster on one machine is needed. For a
  persistent/edge/production node on a host, prefer the `k3s` skill.
---

# kind

kind runs Kubernetes clusters where each "node" is a container. It is a CNCF
project, primarily used for testing Kubernetes itself and for fast, disposable
local dev/CI clusters. Requires a container runtime (Docker or Podman) and the
`kind` and `kubectl` binaries.

## When to use this skill

- Spinning up a throwaway local cluster for development or e2e tests.
- Multi-node topologies (multiple control-plane/worker nodes) on one laptop.
- Loading locally-built images into a cluster without a registry.
- CI pipelines that already have Docker/Podman available.

Use `k3s` instead for a long-lived node on a VM/edge device. Use `podman`
when you only need to build/run images, not orchestrate them.

## Install

```bash
# Go install, or download the release binary for your OS/arch:
go install sigs.k8s.io/kind@v0.23.0
# or
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.23.0/kind-linux-amd64
chmod +x ./kind && sudo mv ./kind /usr/local/bin/kind

kind version
```

## Quick start

```bash
kind create cluster                      # creates cluster named "kind"
kind create cluster --name dev           # named cluster
kubectl cluster-info --context kind-dev  # context is kind-<name>
kind get clusters
kind delete cluster --name dev
```

`kind create cluster` automatically writes/merges the kubeconfig and switches
your current context to `kind-<name>`.

## Pinning the Kubernetes version

Use a node image digest from the matching kind release notes for reproducibility:

```bash
kind create cluster --image kindest/node:v1.30.0
```

## Multi-node clusters (config file)

```yaml
# kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    # expose host ports for an ingress controller
    extraPortMappings:
      - containerPort: 80
        hostPort: 80
      - containerPort: 443
        hostPort: 443
  - role: worker
  - role: worker
```

```bash
kind create cluster --name multi --config kind-config.yaml
```

## Loading local images (the killer feature)

kind nodes can't see your host's image store, so either load or push:

```bash
docker build -t myapp:dev .          # or: podman build -t myapp:dev .
kind load docker-image myapp:dev --name dev
# from a tar:
kind load image-archive myapp.tar --name dev
```

In manifests, set `imagePullPolicy: IfNotPresent` (or `Never`) for loaded
images so the kubelet doesn't try to pull `myapp:dev` from a registry.

## Running on Podman instead of Docker

```bash
export KIND_EXPERIMENTAL_PROVIDER=podman
kind create cluster
```

Rootless Podman needs cgroup v2 and some sysctl/delegation setup — see the
kind "rootless" docs. On rootless, `extraPortMappings` and certain mounts may
behave differently.

## Ingress (common pattern)

Create the cluster with the port mappings shown above, then:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

## Gotchas

- Clusters are **ephemeral** — `kind delete cluster` destroys all state. Don't
  use kind for anything you need to persist.
- LoadBalancer Services have no external IP by default; use `extraPortMappings`
  + NodePort/ingress, or install MetalLB / `cloud-provider-kind`.
- Each node is a container, so cluster resources are bounded by your Docker/
  Podman VM limits (especially on macOS/Windows).
- Image changes won't propagate by rebuilding alone — you must re-run
  `kind load docker-image` and restart the pods.
- Deleting the Docker/Podman engine or pruning containers will silently destroy
  running kind clusters.

## References

- Docs: https://kind.sigs.k8s.io
- Releases (with node image digests): https://github.com/kubernetes-sigs/kind/releases
