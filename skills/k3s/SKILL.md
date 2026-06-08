---
name: k3s
description: >-
  Provision, operate, and tear down lightweight Kubernetes clusters with k3s.
  Use when the user wants a single-binary/edge/IoT/CI Kubernetes cluster,
  mentions k3s, k3sup, or needs a low-footprint production-capable cluster on a
  VM, bare metal, or Raspberry Pi. For throwaway clusters that run *inside*
  Docker/Podman containers, prefer the `kind` skill instead.
---

# k3s

k3s is a CNCF-certified, single-binary Kubernetes distribution (<70 MB, ~512 MB
RAM minimum) from Rancher/SUSE. It bundles containerd, Flannel CNI, CoreDNS,
Traefik ingress, and a local-path storage provisioner, and defaults to SQLite
(swappable for etcd/external DB) instead of a separate etcd. It targets edge,
IoT, CI, and resource-constrained production.

## When to use this skill

- Standing up a real, persistent Kubernetes node on a Linux host/VM/Pi.
- Single-node or multi-node (HA) clusters outside of containers.
- CI runners that need a fast, real cluster without Docker-in-Docker.

If the user wants a disposable cluster running inside containers on their
laptop, use `kind`. If they just need to build/run OCI images, use `podman`.

## Quick start (single node)

```bash
# Install + start the server (systemd service "k3s")
curl -sfL https://get.k3s.io | sh -

# kubeconfig is written here; point kubectl at it
sudo cat /etc/rancher/k3s/k3s.yaml          # contains admin credentials
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml  # or copy to ~/.kube/config

# k3s ships its own kubectl
sudo k3s kubectl get nodes
```

The install script auto-detects the service manager (systemd/openrc). Verify
with `sudo systemctl status k3s` and logs via `sudo journalctl -u k3s -f`.

## Common install options

Pass flags via the `INSTALL_K3S_EXEC` env var or after `sh -s -`:

```bash
# Pin a version (recommended for reproducibility — check channels at
# https://github.com/k3s-io/k3s/releases)
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.30.2+k3s1 sh -

# Disable bundled components you intend to replace
curl -sfL https://get.k3s.io | sh -s - --disable traefik --disable servicelb

# Use a non-default CNI (disable Flannel, then install Calico/Cilium yourself)
curl -sfL https://get.k3s.io | sh -s - --flannel-backend=none --disable-network-policy
```

## Multi-node / HA

```bash
# On the first server, capture the node token
sudo cat /var/lib/rancher/k3s/server/node-token

# Join an AGENT (worker) node
curl -sfL https://get.k3s.io | K3S_URL=https://<server-ip>:6443 \
  K3S_TOKEN=<node-token> sh -

# HA control plane: first server uses embedded etcd via --cluster-init,
# additional servers join with --server
curl -sfL https://get.k3s.io | sh -s - server --cluster-init
curl -sfL https://get.k3s.io | sh -s - server \
  --server https://<first-server-ip>:6443 --token <node-token>
```

HA control planes need an odd number of server nodes (3+) for etcd quorum.

## Operating the cluster

```bash
sudo k3s kubectl get pods -A          # bundled components live in kube-system
sudo k3s crictl ps                    # containerd containers (no docker)
sudo k3s ctr images ls                # containerd image store
```

To use a remote kubeconfig from your workstation, copy `/etc/rancher/k3s/k3s.yaml`
and replace `127.0.0.1` in the `server:` field with the node's reachable IP.

## Teardown

```bash
# Server node
/usr/local/bin/k3s-uninstall.sh
# Agent node
/usr/local/bin/k3s-agent-uninstall.sh
```

## Gotchas

- The kubeconfig at `/etc/rancher/k3s/k3s.yaml` is root-owned (mode 600). Don't
  `chmod` it world-readable on shared hosts; copy it for your user instead.
- Bundled Traefik + ServiceLB (Klipper) claim ports 80/443 — `--disable` them
  if you bring your own ingress/load balancer.
- Default datastore is SQLite (single-server only). For HA you must use
  embedded etcd (`--cluster-init`) or an external DB.
- k3s uses its own containerd, not Docker. Images built with Docker/Podman must
  be pushed to a registry or imported via `k3s ctr images import`.
- On firewalled hosts, open 6443/tcp (API), and for Flannel VXLAN 8472/udp
  between nodes.

## References

- Docs: https://docs.k3s.io
- Releases/channels: https://github.com/k3s-io/k3s/releases
- k3sup (remote bootstrap helper): https://github.com/alexellis/k3sup
