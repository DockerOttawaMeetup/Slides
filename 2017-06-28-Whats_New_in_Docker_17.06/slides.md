layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# ✨ What's new in Docker 17.06? ✨

_Dave Henderson, June 2017_

---

## Docker's release schedule 🗓

<img src="https://i0.wp.com/blog.docker.com/wp-content/uploads/lifecycle.png" width="95%">

- This is the first major Long-Term Support (LTS) release of Docker CE/EE (17.03 was really just 13.2)
- 17.06 is the first version of Docker built entirely on the Moby Project (http://mobyproject.org/).

---

## Old new Features

- **Lots** of new stuff for those on the `stable` channel
- From [17.04](https://github.com/moby/moby/releases/tag/v17.04.0-ce):
  - Hosts can now join swarm with `--availability=drain` (allows "manager-only from birth") [#24993](https://github.com/moby/moby/pull/24993)
  - Topology-aware scheduling (`--placement-pref`), uses node/engine labels, e.g. for spreading replicas evenly across datacenters/racks [#30725](https://github.com/moby/moby/pull/30725)
  - Automatic service rollback on failure (`--update-failure-action=rollback`) [#31108](https://github.com/moby/moby/pull/31108)
  - Add `--read-only` for service create and service update - allows immutable containers in Swarm too [#30162](https://github.com/moby/moby/pull/30162)
- From [17.05](https://github.com/moby/moby/releases/tag/v17.05.0-ce):
  - Multi-stage builds [#31257](https://github.com/moby/moby/pull/31257) [#32063](https://github.com/moby/moby/pull/32063)
  - Allow using build-time args (`ARG`) in `FROM` [#31352](https://github.com/moby/moby/pull/31352)
  - [Logging driver plugins](https://docs.docker.com/engine/admin/logging/plugins/) (provide a log driver via a container) [#28403](https://github.com/moby/moby/pull/28403)
  - Fix UDP traffic in containers not working after the container is restarted [#32505](https://github.com/moby/moby/pull/32505)
  - Synchronous service create and service update (`--detach=true/false`), adds progress bar and failure indication [#31144](https://github.com/moby/moby/pull/31144)
  - New `--start-period` option for `HEALTHCHECK` [#28938]((https://github.com/moby/moby/pull/28938)
  - Move `docker service logs` out of experimental [#32462]((https://github.com/moby/moby/pull/32462)

---

### New `docker build --iidfile` option [#32406](https://github.com/moby/moby/pull/32406)

Allows specifying a location where to save the resulting image ID. Useful for situations where you want to avoid tagging with a name, like some CI builds.

```Makefile
image.id: Dockerfile
    docker build --iidfile=image.id .
other-target: image.id
    do-stuff-with <image.id
```

### Allow specifying any remote ref in git checkout URLs [#32502](https://github.com/moby/moby/pull/32502)

Can now build an image direct from a GitHub PR (there must be a `Dockerfile` in the root fo the repo in that branch).

```console
$ docker build git://github.com/hairyhenderson/dockerfiles#pull/10/head
...
```

---

### More commands have `--format` option [#31557](https://github.com/moby/moby/pull/31557) [#30962](https://github.com/moby/moby/pull/30962) [#31482](https://github.com/moby/moby/pull/31482)

```console
$ docker stack ls --format '{{.Name}}'
...
$ docker history --format '{{.CreatedBy}}' hairyhenderson/gomplate
...
$ docker system df --format '{{.Active}} {{.Type}} - {{.Reclaimable}}'
```

_Pro Tip: use `{{ . | json }}` to see available keys_

---

### New stack file options [docker/cli#73](https://github.com/docker/cli/pull/73) [docker/cli#35](https://github.com/docker/cli/pull/35) [#32059](https://github.com/moby/moby/pull/32059)

```yaml
version: '3.3'
services:
  myservice:
    ...
    read_only: true
    ...
    deploy:
      placement:
        preferences:
          - spread: node.labels.rack
    ...
    dns:
      - 8.8.8.8
    dns_search:
      - foo.com
```

-  `read_only` - enables immutable containers in services deployed by `docker stack deploy`.
- placement preferences - just like `--placement-pref` in `docker service create`, for topology-aware scheduling
- `dns` and `dns_search` - allows specifying nameservers and search domains in stack files
  - _**by Ottawa's own [Ben Boudreau](https://github.com/seriousben)!**_

---

### Display Swarm cluster and node TLS information [docker/cli#44](https://github.com/docker/cli/pull/44)

```console
$ docker system info
...
Swarm: active
...
 CA Configuration:
  Expiry Duration: 3 months
  Force Rotate: 0
 Root Rotation In Progress: false
...
```

_Note: `docker system info` is equivalent to the sort-of-deprecated `docker info`_

### New `docker swarm ca` subcommand [docker/cli#48](https://github.com/docker/cli/pull/48)

Allows managing a swarm CA

```console
$ docker swarm ca --rotate --cert-expiry 4h
desired root digest: sha256:237f91194c9e10b3f569b9e883985e0fffee24761ba43b8cb882b31d095f828d 
  rotated TLS certificates:  1/1 nodes
  rotated CA certificates:   1/1 nodes
...
$ docker system info --format '{{.Swarm.Cluster.Spec.CAConfig.NodeCertExpiry}}'
4h0m0s
```

---

### Long-format options supported on `--network` and `--network-add` [docker/cli#62](https://github.com/docker/cli/pull/62) [#33130](https://github.com/moby/moby/pull/33130)

Rare case of syntax being supported in stack file first - adds long-form (CSVish) options for networks with `docker service create`/`update`.

```console
$ docker service create --network name=docknet,alias=web1,driver-opt=field1=value1
...
```

### More accurate `docker stats` memory output [docker/cli#80](https://github.com/docker/cli/pull/80)

Used to show used + buffered/cached, now shows _just_ used

```console
$ docker stats --no-stream --format '{{.MemUsage}}'
1.785MiB / 1.952GiB
```

---

### Add multiline processing to the AWS CloudWatch logs driver [#30891](https://github.com/moby/moby/pull/30891)

Two new log options for the `awslogs` log driver:
- `awslogs-datetime-format`: start pattern for multi-line messages that start with timestamps
- `awslogs-multiline-pattern`: start pattern in regex format

_This is useful for keeping multi-line stack traces in one log event in AWS!_

```console
$ docker run --log-driver=awslogs --log-opt awslogs-group=test --log-opt awslogs-datetime-format='%Y-%m-%d' ...
```

### Add option to auto-configure blkdev for devmapper [#31104](https://github.com/moby/moby/pull/31104)

Instead of forcing users to manually configure a block device to use with devmapper, this gives the user the option to let the devmapper driver configure a device for them.

---

### Swarm-mode services can now attach to node-local networks [#32981](https://github.com/moby/moby/pull/32981)

Previously, swarm services could only attach to "swarm-scoped" networks, so you couldn't attach to the `host` network (i.e. the host's local interface).

Now, you can attach to `host`, `bridge`, and other node-local networks.

```console
$ docker service create --name myservice --network host busybox ifconfig docker0
...
$ docker service logs -f myservice
myservice.1.6lajq1vlefdd@moby    | docker0   Link encap:Ethernet  HWaddr 02:42:C3:59:7A:BD  
...
```

### Isolate Swarm Control-plane traffic from Application data traffic using --data-path-addr [#32717](https://github.com/moby/moby/pull/32717)

Allows using separate networks for cluster control messages and application data (overlay networks, etc).

```console
$ docker swarm init --advertise-addr 10.0.0.1 --datapath-addr 192.168.0.1
```

---

### More prometheus-format metrics [#32792](https://github.com/moby/moby/pull/32792) [docker/swarmkit#2157](https://github.com/docker/swarmkit/pull/2157)

More metrics:
- `swarm_node_info`
- `swarm_node_info{node_id,swarm_id}`
- `swarm_manager_nodes{state="disconnected|down|ready|unknown"}`
- `swarm_node_manager 0/1`
- `builder_builds_failed_total{reason}`
- `builder_builds_triggered_total`
- `engine_daemon_engine_info{architecture,commit,daemon_id,graphdriver,kernel,os,os_type,version}`

### New metric plugins support [#32874](https://github.com/moby/moby/pull/32874)

Allows a new type of plugin that can be used to scrape or proxy access to the Docker Engine's prometheus-format metrics.

```console
$ docker plugin install --grant-all-permissions cpuguy83/docker-metrics-plugin-test:latest
$ docker run --rm --net host buildpack-deps:curl curl localhost:19393/metrics
```

---

### Cluster events now broadcst in Docker event stream [#32421](https://github.com/moby/moby/pull/32421)

```console
$ docker events --filter 'scope=swarm'
2017-06-27T00:08:38.573347085-04:00 service create le61se43rg2e5sim4pdoxpa4a (name=myservice)
2017-06-27T00:08:38.574232652-04:00 service update le61se43rg2e5sim4pdoxpa4a (name=myservice)
2017-06-27T00:08:54.464902006-04:00 service remove le61se43rg2e5sim4pdoxpa4a (name=myservice)
```

### Allow specifying a secret location within the container [#32571](https://github.com/moby/moby/pull/32571)

Allows storing secrets in locations outside of `/run/secrets/`.

```console
$ docker service create --secret source=foo,target=/etc/foo ...
```

### Initial support for secrets on Windows [#32208](https://github.com/moby/moby/pull/32208)

- Not stored in RAM (lack of RAM filesystem supoprt in Windows)
- Available inside the container in `C:\ProgramData\Docker\internal\secrets`
- Not quite as configurable (UID/GUID/mode not settable, only accessible by `Administrator` and `SYSTEM`)

---

### New config support [#32336](https://github.com/moby/moby/pull/32336), [docker/cli#45](https://github.com/docker/cli/pull/45), [#33169](https://github.com/moby/moby/pull/33169)

Add support for services to carry arbitrary config objects.

Like secrets, except value is readable through `docker config inspect` command. And not encrypted in the container.

```console
$ echo 'foo=bar' | docker config create fooconf -
n8ovojzfx4x5ntwwnu544c2kp
$ docker config inspect fooconf
[
    {
        "ID": "n8ovojzfx4x5ntwwnu544c2kp",
        ...
        "Spec": {
          ...
            "Data": "Zm9vPWJhcgo="
...
$ echo 'Zm9vPWJhcgo=' | base64 -d
foo=bar
```

---

## Removed/Deprecated Features 👋

### Disable legacy registry (v1) by default [#33629](https://github.com/moby/moby/pull/33629)

Docker 17.06 _by default_ disables communication with legacy (v1) registries. Can still use `--disable-legacy-registry=false` on `dockerd` to force re-enabling this.

To be _totally_ removed in Docker 17.12.

### Remove deprecated `--email` flag from `docker login` [docker/cli#143](https://github.com/docker/cli/pull/143)

```console
$ docker login --email
unknown flag: --email
```

---

## Digging Deeper 🕳

- 17.04 changelog: https://github.com/moby/moby/releases/tag/v17.04.0-ce
- 17.05 changelog: https://github.com/moby/moby/releases/tag/v17.05.0-ce
- 17.06 changelog: https://github.com/docker/docker-ce/releases/tag/v17.06.0-ce _(note the repo change)_

---

## Thank you!

### Any Questions?
