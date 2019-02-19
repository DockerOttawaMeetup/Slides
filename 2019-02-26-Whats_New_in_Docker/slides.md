layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: center, middle

# ✨ What's new in Docker? ✨

_also: what's not so new..._

_Dave Henderson, February 2019_

---

## What happened in 2018? 🗓

- 4 monthly _edge_ releases (18.01, 18.02, 18.04, 18.05)
- 3 quarterly _stable_ releases (18.03, 18.06, 18.09)
- Switched to ~6 month release cycle, starting with 18.09 _(with 7mo support for CE)_
  - _synchronized releases between CE & EE, on all platforms_

---

## Features from early 2018... 

- _Experimental_ `docker manifest` command [docker/cli#138](https://github.com/docker/cli/pull/138) _(v18.02)_
  - See [Multi Arch Dockerfiles](../2018-02-28-Multi_Arch_Dockerfiles/slides.md) talk from February 2018
  - _still_ experimental...
- Allow `Dockerfile` to be outside of build-context [docker/cli#886](https://github.com/docker/cli/pull/886) _(v18.03)_
- `docker trust` commands graduate out of experimental [docker/cli#934](https://github.com/docker/cli/pull/934) / [docker/cli#935](https://github.com/docker/cli/pull/935) / [docker/cli#944](https://github.com/docker/cli/pull/944) _(v18.03)_
  - cryptographic signing/verification of images
- SCTP port mapping support [docker/cli#278](https://github.com/docker/cli/pull/278) / [docker/swarmkit#2298](https://github.com/docker/swarmkit/pull/2298) _(v18.03)_
  - `sctp` services can now be published from stand-alone containers and Swarm services

---

## Notable features in 18.06

### BuildKit! 🏗

- New _experimental_ builder backend based on BuildKit [moby/moby#37151](https://github.com/moby/moby/pull/37151) / [docker/cli#1111](https://github.com/docker/cli/pull/1111)
  - _More on this later_

### Kubernetes <img src="images/kubernetes.png" height="30em">

- Docker Compose/Stacks on Kubernetes commands graduate from experimental. [docker/cli#1068](https://github.com/docker/cli/pull/1068) [docker/cli#899](https://github.com/docker/cli/pull/899)
  - `docker stack deploy --orchestrator=kubernetes ...`
  - [docker/compose-on-kubernetes](https://github.com/docker/compose-on-kubernetes) open-sourced at DockerCon EU 2018, can be installed on any Kubernetes cluster!

---

## Notable features in 18.09

### Deprecations are features! 💀

- Deprecated `devicemapper` and legacy `overlay` (v1) storage drivers [docker/cli#1455](https://github.com/docker/cli/pull/1455) / [docker/cli#1424](https://github.com/docker/cli/pull/1424) / [docker/cli#1425](https://github.com/docker/cli/pull/1425)
- Removed support for TLS < 1.2 [moby/moby#37660](https://github.com/moby/moby/pull/37660)
- Removed Ubuntu 14.04 and Debian 8 as supported platforms [docker/docker-ce-packaging#255](https://github.com/docker/docker-ce-packaging/pull/255) / [docker/docker-ce-packaging#254](https://github.com/docker/docker-ce-packaging/pull/254)

_(See https://docs.docker.com/engine/deprecated/ for removal dates)_

---

## Notable features in 18.09

### Connection to Docker Engine over SSH 🔐

- Add support for remote connections using SSH [docker/cli#1014](https://github.com/docker/cli/pull/1014)
  - `docker -H ssh://me@my-server run hello-world`
  - or `export DOCKER_HOST=ssh://me@server`
  - just shells out to `ssh` so `~/.ssh/config` is honoured

---

## Notable features in 18.09

### More BuildKit! 🏗

- BuildKit graduates from experimental [moby/moby#37593](https://github.com/moby/moby/pull/37593) / [moby/moby#37686](https://github.com/moby/moby/pull/37686) / [moby/moby#37692](https://github.com/moby/moby/pull/37692) / [docker/cli#1303](https://github.com/docker/cli/pull/1303) / [docker/cli#1275](https://github.com/docker/cli/pull/1275)
  - Still not default, but can now also be configured with an option in `daemon.json`
- Support for build-time secrets with `--secret` flag [docker/cli#1288](https://github.com/docker/cli/pull/1288)
- SSH agent socket forwarder (`docker build --ssh`) [docker/cli#1438](https://github.com/docker/cli/pull/1438) / [docker/cli#1419](https://github.com/docker/cli/pull/1419)
- New `builder prune` subcommand to prune BuildKit build cache [docker/cli#1295](https://github.com/docker/cli/pull/1295) [docker/cli#1334](https://github.com/docker/cli/pull/1334)

---

## Notable features in 18.09

### Logging 📃

- Add `local` log driver [moby/moby#37092](https://github.com/moby/moby/pull/37092)
  - Much more efficient compressed binary format (protobuf)
  - Will eventually replace `jsonlogfile`
    - (once dual-driver logging is implemented!)

---

## Notable features in 18.09

### Packaging 📦

- Removed Ubuntu 14.04 "Trusty Tahr" and Debian 8 "Jessie" as supported platforms [docker/docker-ce-packaging#255](https://github.com/docker/docker-ce-packaging/pull/255) / [docker/docker-ce-packaging#254](https://github.com/docker/docker-ce-packaging/pull/254)
- Split `engine`, `cli`, and `containerd` to separate packages, and run `containerd` as a separate `systemd` service docker-ce-packaging#131, docker-ce-packaging#158
- Remove `-ce` / `-ee` suffix from version string [docker/docker-ce-packaging#206](https://github.com/docker/docker-ce-packaging/pull/206)

---

## Demos 🎭

- Connection to Docker Engine over SSH
- compose-on-kubernetes

---

## Digging Deeper 🕳

- `18.09.2` changelog: https://github.com/docker/docker-ce/releases/tag/v18.09.2 (includes release notes for `.1` and `.0`)
- compose-on-kubernetes: https://github.com/docker/compose-on-kubernetes

---

## Thank you! 🙇‍

### Any Questions?
