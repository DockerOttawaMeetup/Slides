layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# ✨ What's new in Docker? ✨

_Dave Henderson, October 2017_

---

## Docker's release schedule 🗓

<img src="https://i0.wp.com/blog.docker.com/wp-content/uploads/lifecycle.png" width="95%">

---

## Notable features in 17.07

_17.09 includes all features from the 17.07 edge release, including these notable ones:_

- Passwords can now be passed using `STDIN` using the new `--password-stdin` flag on `docker login` [docker/cli#271](https://github.com/docker/cli/issues/271)
  - _(related)_ Print a warning if docker login `--password` is used, and recommend `--password-stdin` [docker/cli#270](https://github.com/docker/cli/issues/270)
- Allow stopping of paused container [moby/moby#34027](https://github.com/moby/moby/issues/34027)
- Add quota support for the overlay2 storage driver [moby/moby#32977](https://github.com/moby/moby/issues/32977)
  - _when the backing filesystem is `xfs`_
- Initial support for plugable secret backends [moby/moby#34157](https://github.com/moby/moby/issues/34157) [moby/moby#34123](https://github.com/moby/moby/issues/34123)

---

## Notable features in 17.09

- Add `--chown` flag to `ADD`/`COPY` commands in `Dockerfile` [moby/moby#34263](https://github.com/moby/moby/issues/34263)
- Support for Compose v3.4 format
  - Allow extension fields in the v3.4 version of the compose format [docker/cli#452](https://github.com/docker/cli/issues/452)
      - _allows `x-*` fields for integration with 3rd-party tooling_ ([example](https://medium.com/@kinghuang/docker-compose-anchors-aliases-extensions-a1e4105d70bd))
  - Make compose file allow to specify names for non-external volume [docker/cli#306](https://github.com/docker/cli/issues/306)
      - _Allows use of templates in volume names - fixes [moby/moby#30770](https://github.com/moby/moby/issues/30770)_
  - Support `start_period` for `healthcheck` in Docker Compose [docker/cli#475](https://github.com/docker/cli/issues/475)
  - Add support for update order in compose deployments [docker/cli#360](https://github.com/docker/cli/issues/360)
      - _set to `start-first` to reduce downtime_
- Graphdriver: promote `overlay2` over `aufs` [moby/moby#34430](https://github.com/moby/moby/issues/34430)
- Add `docker service rollback` subcommand [docker/cli#205](https://github.com/docker/cli/issues/205)
    - _allows explicitly rolling back to previous service configuration_
- Fix managers failing to join if the gRPC snapshot is larger than 4MB [docker/swarmkit#2375](https://github.com/docker/swarmkit/issues/2375)
- Fix services failing to deploy on ARM nodes [moby/moby#34021](https://github.com/moby/moby/issues/34021)

---

### Notable features in 17.10

_These changes won't be in the `stable` release until 17.12_

- Use non-detached mode as default for `docker service` commands [docker/cli#525](https://github.com/docker/cli/issues/525)
- Add experimental `docker trust`: `view`, `revoke`, `sign` subcommands [docker/cli#472](https://github.com/docker/cli/issues/472)
- Add support for `.Node.Hostname` templating in swarm services [moby/moby#34686](https://github.com/moby/moby/issues/34686)

---

## About the Kubernetes support...

- (private?) beta coming soon (sign up at https://docker.com/kubernetes)
- initial focus is on Docker EE, and Docker for Mac/Windows - Docker CE Linux
  support on roadmap (post 2018 Q1)
- all EE features (DTR, private registry, scanning, etc) will work with both Swarmkit and K8s
- shipping "vanilla" k8s from CNCF - no wrapping involved (using LinuxKit)
- "opinionated" stack - "batteries included" (networking, storage) - and _eventually_ swappable
- _No plans to drop Swarm support_

---

## Digging Deeper 🕳

- 17.07 changelog: https://github.com/docker/docker-ce/releases/tag/v17.07.0-ce
- 17.09 changelog: https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce
- 17.10 changelog: https://github.com/docker/docker-ce/releases/tag/v17.10.0-ce

---

## Thank you!

### Any Questions?
