---
name: podman
description: >-
  Build, run, and manage OCI containers, images, pods, and volumes with Podman,
  a daemonless, rootless-capable Docker-compatible engine. Use when building or
  running containers, writing Containerfiles/Dockerfiles, working with pods or
  Quadlet/systemd units, or migrating from Docker. For orchestrating multi-node
  Kubernetes clusters, use the `kind` or `k3s` skills instead.
---

# podman

Podman is a daemonless container engine with a Docker-compatible CLI. It runs
containers as child processes of the invoking user (no central daemon), supports
rootless mode by default, and introduces the concept of **pods** (groups of
containers sharing a network namespace, like a Kubernetes pod).

## When to use this skill

- Building images from a Containerfile/Dockerfile.
- Running/managing containers, images, volumes, networks.
- Grouping containers into pods, or generating Kubernetes YAML from them.
- Running containers as systemd services via Quadlet.
- Migrating Docker workflows (most `docker` commands map 1:1).

Use `kind`/`k3s` when you need a full Kubernetes control plane. Podman can
*back* a kind cluster (see the `kind` skill) but does not orchestrate on its own.

## Docker compatibility

Almost every `docker` subcommand works as `podman`:

```bash
podman pull docker.io/library/alpine:3.20
podman run --rm -it alpine:3.20 sh
podman ps -a
podman images
```

For drop-in compatibility you can `alias docker=podman`, or install the
`podman-docker` package which provides a `docker` shim and a socket emulating
the Docker API (`podman system service`). Note: Podman requires **fully
qualified image names** or configured registries — `docker.io/library/...`
rather than bare `alpine` on strict configs.

## Building images

```bash
# Containerfile is preferred; Dockerfile also works
podman build -t myapp:dev -f Containerfile .
podman build -t myapp:dev .            # auto-detects Containerfile/Dockerfile

# Inspect / tag / push
podman images
podman tag myapp:dev registry.example.com/team/myapp:dev
podman login registry.example.com
podman push registry.example.com/team/myapp:dev
```

Podman uses Buildah under the hood. For advanced multi-arch builds use
`podman build --platform linux/amd64,linux/arm64 --manifest myapp:multi`.

## Rootless mode

Podman runs rootless by default for non-root users (uses user namespaces +
`subuid`/`subgid` maps from `/etc/subuid` and `/etc/subgid`).

```bash
podman unshare cat /proc/self/uid_map   # inspect the UID mapping
podman info | grep -i rootless
```

Rootless caveats:
- Can't bind host ports < 1024 unless `net.ipv4.ip_unprivileged_port_start` is
  lowered or you run rootful (`sudo podman`).
- Uses slirp4netns/pasta or rootless netavark for networking; performance and
  source-IP behavior differ from rootful.
- Storage lives under `~/.local/share/containers` instead of `/var/lib/containers`.

## Pods

```bash
podman pod create --name web -p 8080:80
podman run -d --pod web nginx:alpine
podman run -d --pod web myapp:dev
podman pod ps
podman pod stop web && podman pod rm web
```

## Generate / play Kubernetes YAML

```bash
# Export a running pod/container to Kubernetes manifests
podman generate kube web > web.yaml      # (newer: podman kube generate)

# Recreate workloads from a manifest
podman play kube web.yaml                 # (newer: podman kube play)
```

This is a convenient bridge to `kind`/`k3s`: prototype locally, then apply the
generated YAML to a real cluster.

## Running as systemd services (Quadlet)

Prefer **Quadlet** (`.container`, `.pod`, `.volume`, `.network` unit files in
`/etc/containers/systemd` or `~/.config/containers/systemd`) over the deprecated
`podman generate systemd`:

```ini
# ~/.config/containers/systemd/myapp.container
[Container]
Image=registry.example.com/team/myapp:dev
PublishPort=8080:80

[Install]
WantedBy=default.target
```

```bash
systemctl --user daemon-reload
systemctl --user start myapp.service
```

## podman-compose

For Compose files, use `podman compose` (delegates to an installed compose
provider) or the `podman-compose` Python tool:

```bash
podman compose up -d
```

## Gotchas

- Daemonless: there is no background daemon to "restart". Container lifecycle
  ties to the invoking process; use systemd/Quadlet for persistence across
  reboots.
- Short image names can fail or prompt for a registry — qualify them
  (`docker.io/library/...`) or set `unqualified-search-registries` in
  `/etc/containers/registries.conf`.
- Rootless containers can't reach privileged host ports and have a separate
  storage root from rootful — `sudo podman ps` and `podman ps` show different
  containers.
- `docker generate systemd` is deprecated; use Quadlet for new work.
- macOS/Windows run Podman inside a managed VM via `podman machine init` /
  `podman machine start` — remember to start it first.

## References

- Docs: https://docs.podman.io
- Quadlet: https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html
- Rootless tutorial: https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md
