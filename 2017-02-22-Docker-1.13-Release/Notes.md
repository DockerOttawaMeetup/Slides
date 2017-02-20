layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# What's new in Docker 1.13

_Dave Henderson, February 2017_

???

_Adapted from recent talk given by Victor Vieux_

---

## Restructured CLI commands

_[#26025](https://github.com/docker/docker/issues/26025)_

The `docker` command has lots of newer "management" subcommands (`volume`, `network`, `service`, `node`, etc...) - 2 new subcommands added: `image` and `container`, to contain all of the previously-top-level commands.

- `docker images` -> `docker image list`
- `docker build` -> `docker image build`
- `docker run` -> `docker container run`
- `docker create` -> `docker container create`
- _you get the idea..._

- fully backward compatible, legacy commands still shown in `--help`.
  - Use `DOCKER_HIDE_LEGACY_COMMANDS=1` to hide legacy commands.
- will be deprecated "Maybe in a year or two"
- certain commands (`run`, `build`) remain as permanent aliases and probably won't ever be deprecated

---

## Experimental no longer separate binary

_[#27223](https://github.com/docker/docker/issues/27223)_

Now, `dockerd` has a `--experimental` argument.

Makes enabling experimental mode much simpler.

Turned on by default in:
- Docker for Mac/Windows beta channels
- Docker for AWS/Azure stable (for now) and beta channels

Note: _Experimental_ does not mean _unstable_. These are just features where the design/API may change in future releases.

---

## CLI backward compatibility

_[#27745](https://github.com/docker/docker/pull/27745)_

The `docker` CLI can now talk to older daemons.

No more `Client and server don't have the same version` error!

```console
$ docker version
Client:
 Version:      1.13.0-dev
 API version:  1.24 (downgraded from 1.25)
...

Server:
 Version:      1.12.2
 API version:  1.24
...

$ docker create --stop-timeout 5 test
Server only supports API version stop timeout, "1.24" support requires at least 1.25
```

---

## Default swarm mode encryption at rest

_[#27967](https://github.com/docker/docker/pull/27967)_

- All of the data stored in swarm mode is now encrypted at rest (for swarms created with 1.13)
- Adds new "manager autolocking" feature.
  - When enabled, swarm manager nodes must be unlocked with a key after restart (with `docker swarm unlock` command).

```console
docker swarm init --autolock
docker swarm update --autolock
docker swarm unlock
docker swarm unlock-key --rotate
```

---

## `docker plugin` command out of experimental

_[#28226](https://github.com/docker/docker/pull/28226)_

This is the "v2" plugin API: allows installing/managing network & volume plugins via a `docker plugin` command.

e.g.
```console
docker plugin create vieux/sshfs /path/to/rootfs
docker plugin enable vieux/sshfs
docker plugin install vieux/sshfs
docker plugin set vieux/sshfs DEBUG=1
```

- 1.13.1 added `docker plugin upgrade` for _almost_ zero-downtime plugin upgrades

---

## Use docker-compose file to deploy to swarm

_[#27998](https://github.com/docker/docker/pull/27998)_

New `--compose-file`/`-c` argument to `docker stack deploy`!

```console
docker stack deploy --compose-file=docker-compose.yml foo
docker stack list
docker stack rm foo
```

- New v3 compose format
  - Removed the non portable options (build, volume-from, etc...) from v2 format
  - Added Swarm specific options (replicas, mode, etc...)
- Docker 1.13.1 adds v3.1 compose format
  - support for secrets

---

## Data management commands

_[#26108](https://github.com/docker/docker/pull/26108)_

- `docker system df` - show Docker's disk usage
- `docker system prune` - remove unused data
- Also `prune` subcommands for other resources (`container`, `image`, `volume`, `network`).
- Will not delete all unused images by default (use `--all`/`-a` arg for that)

---

## Swarm mode secret management

_[#27794](https://github.com/docker/docker/pull/27794)_

Allows secrets to be securely provided to swarm services. _Not yet available to standalone containers._

```console
$ echo "password1" | docker secret create foo -
$ docker service create --secret foo --name myservice alpine cat /run/secrets/foo
$ docker service logs myservice
myservice.1.w3ckid9mq0i9@moby    | password1
```

- secret is encrypted end-to-end
- replicated by swarm managers, and only made available to hosts & services that need the secret
- flushed from memory after container stops running
- secret available to container in `/run/secrets` (in-memory temporary filesystem)
- name can be customized with more advanced mount-style syntax:
  ```
  $ docker service create --name myservice \
    --secret source=foo,target=password ...
  ```
- secret can't be deleted while in use
  - secret rotation supported with `--secret-add`/`--secret-rm` options on `docker service update` - simply use versioned secret names

---

## Swarm mode templating for options

_[#28025](https://github.com/docker/docker/pull/28025)_

- used on `service create`
- uses Go's [text/template](http://golange.org/pkg/text/template/) syntax
- supported by `--hostname`, `--mount`, `--env`

```console
$ docker service create --name foo --hostname "{{.Service.Name}}-{{.Task.Slot}}" --replicas 3 alpine hostname
$ docker service logs foo
foo.1.tbtj4wqnmou1@moby    | foo-1
foo.2.kt926okv9s6q@moby    | foo-2
foo.3.47v0spa7ire8@moby    | foo-3
```

---

## Swarm mode templating for options (continued)

- _real_ handy for mounting the same volume to the same task slot
- supported in compose file too, though volume syntax is _awkward_:
    ```yaml
    version: '3.1'
    services:
      redis:
        image: redis
        volumes:
          - redisVol:/data
    volumes:
      redisVol:
        external:
          name: '{{index .Service.Labels "com.docker.stack.namespace"}}_redisVol-{{.Task.Slot}}'
    ```
  - I logged a bug about this: _[#30770](https://github.com/docker/docker/issues/30770)_

---

## Swarm mode health-aware load-balancing and DNS resolution

_[#27279](https://github.com/docker/docker/pull/27279)_

When using the `HEALTHCHECK` Dockerfile instruction (or `--health-cmd` on `docker run`), swarm mode LB will now only start sending traffic once healthcheck is passing.

---

## Swarm mode support in Windows

_[#27838](https://github.com/docker/docker/pull/27838)_

_But overlay networking not yet supported, so may not be that useful yet..._

---

## (experimental) `build --squash`

_[#22641](https://github.com/docker/docker/pull/22641)_

Allows built images to be flattened to the `FROM`, saving space in final image.

- does not destroy any images or layers (on the build machine)
- preserves build cache (for fast re-builds)
- be careful using this feature for injecting secrets with build args, because history is preserved

```console
$ docker build --squash -t hairyhenderson/figlet .
...
$ docker history hairyhenderson/figlet
IMAGE               CREATED             CREATED BY                                      SIZE                COMMENT
68675befc59b        3 seconds ago                                                       657 kB              merge sha256:2046cd62da11a67fe832c201c14e8b29ca57b92e5241b7ffd2509a12bd948d66 to sha256:88e169ea8f46ff0d0df784b1b254a15ecfaf045aee1856dca1ec242fdd231ddd
<missing>           53 seconds ago      /bin/sh -c #(nop)  ENTRYPOINT ["figlet"]        0 B
<missing>           54 seconds ago      /bin/sh -c apk update   && apk add figlet ...   0 B
<missing>           7 weeks ago         /bin/sh -c #(nop) ADD file:92ab746eb22dd3e...   3.98 MB
```

---

## (experimental) `docker service logs`

_[#28089](https://github.com/docker/docker/pull/28089)_

Shows aggregated logs from all instances of a swarm service.

---

## Zombie-reaping init support with `--init`

_[#26061](https://github.com/docker/docker/pull/26061)_

- Helps in some situations where containers hang due to signal-handling restrictions on processes being PID 1.
- Uses a very small `init` replacement called `tini` (but can be overridden with `--init-path`)

```console
$ docker run alpine ps
PID   USER     TIME   COMMAND
    1 root       0:00 ps
$ docker run --init alpine ps
PID   USER     TIME   COMMAND
    1 root       0:00 /dev/init -- ps
    6 root       0:00 ps
```

---

## Build-time network attachment for `RUN`s

_[#27702](https://github.com/docker/docker/pull/27702)_

- `docker build --network`
- allows `RUN` instructions in `Dockerfile` to have access to other networks.
- can attach to a running container's own network,
  or direct to host network,
  or a named network,
  or no network (prevent commands from accessing network)
- can be used for retrieving secrets at build time in a secure manner

---

## Deprecations

- `docker daemon` command
  - _[#26834](https://github.com/docker/docker/pull/26834)_
  - use `dockerd` instead (introduced in 1.12)
- `MAINTAINER` instruction
  - _[#25466](https://github.com/docker/docker/pull/25466)_
  - use `LABEL maintainer "dhenderson@gmail.com"` instead

---

## Thank you!
