layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# Stop Running Containers As `root`!

No really, it's a terrible idea!

_Dave Henderson, May 2018_

---

## Thank you!

### Any Questions?

---

## Why not?

Common reasons given to run as `root`:

- _The container is isolated, so there's no harm!_
- _I need to bind to port 80 (or some other port < 1024)_
- _I want to install software..._
- _My application just expects to be root!_

We'll deal with these one by one, but first, how to avoid this altogether?

---

## Running containers without `root`

With `docker run`:
```console
$ docker run -u nobody -it alpine
~ $ id
uid=65534(nobody) gid=65534(nobody)
```

--

With `docker exec`:
```console
$ docker exec -u nobody -it mycontainer id
uid=65534(nobody) gid=65534(nobody)
```

---

## Running containers without `root`

With `docker-compose`:
```yaml
version: '3'

services:
  foo:
    image: alpine
    user: nobody
    command: id
```
```console
$ docker-compose up
...
foo_1  | uid=65534(nobody) gid=65534(nobody)
demo_foo_1 exited with code 0
```

_Same syntax for Docker Swarm stacks..._


---

## Building images to run as non-`root`

Use the [`USER`](https://docs.docker.com/engine/reference/builder/#user) instruction:

```dockerfile
FROM node

RUN mkdir -p /app
COPY . /app
RUN chown -R nobody /app

WORKDIR /app

# Not always necessary - can just use the UID
RUN groupadd -g 64000 fakeuser && \
    useradd -r -u 64000 -g fakeuser fakeuser
USER fakeuser

CMD ["node", "."]
```

---

## Building `FROM scratch` images as non-`root`

```dockerfile
FROM ubuntu:latest
RUN useradd -u 10001 scratchuser

FROM scratch
COPY dosomething /dosomething
COPY --from=0 /etc/passwd /etc/passwd
USER scratchuser
ENTRYPOINT ["/dosomething"]
```

_Or..._

```dockerfile
FROM scratch

COPY dosomething /dosomething

USER 10001:10001

ENTRYPOINT ["/dosomething"]
```

- _pay attention to the GID!_
- _only works if your application doesn't care about usernames (doesn't use the `getuid(2)` function)_

---

## I want to run as `root` because...

### I'm in a container, what's the harm?

--

- `root` in a container is `root` on the host _(usually mitigated by_ namespaces _like `pid`, `net`, and `mnt`)_
  - you wouldn't run a host service as `root`, _would you?_

--

- bind-mounted volumes can be accessed as root:
  ```console
  $ docker run -v /etc:/e --rm alpine touch /e/hi_there
  $ ls -alh /etc/hi_there
  -rw-r--r-- 1 root 0 May 28 19:54 /etc/hi_there
  ```
  - _partly mitigated by read-only host filesystems like with LinuxKit_

--

- _container escape_ vulnerabilities have been found in the past
  - will probably be found in the future
  - exploits require running as `root`

--

Best way to stay safe? Don't run as `root`!

---

## I want to run as `root` because...

### I need to bind to port 80 (or some other port < 1024)

--

- _Almost always_ avoidable
  - configure service to listen to a higher port
  - if it must appear as `80`, publish the port as `80`
  ```console
  $ docker run -p 80:8080 my/server
  ```

--

- If it's not avoidable, and you have a new-enough kernel (4.11+):
  ```console
  $ docker run --sysctl=net.ipv4.ip_unprivileged_port_start=0 ...
  ```

--

- Better yet, use user namespace remapping! _(more on this later)_

---

## I want to run as `root` because...

### I want to install software...

--

This is an anti-pattern! Don't treat your containers like VMs...

Instead, only install software in the `Dockerfile`, and use the `USER` instruction:
```dockerfile
FROM alpine

RUN apk add --no-cache figlet

USER 65000:65000

ENTRYPOINT ["figlet"]
```

- may need to `chown` some files before the `USER` instruction
- when starting `FROM` an image that already has a `USER` instruction, use `USER root` to become root, then `USER <non-root user>` afterwards

---

## I want to run as `root` because...

### My application expects to be `root`!

--

- Fix your application...
  - If it expects to be `root`, this probably violates the [_principle of least privilege_](https://en.wikipedia.org/wiki/Principle_of_least_privilege), so consider this a bug!
  - If it's not your application, or you can't configure or alter it...

--

- Use user namespace remapping!

---

## Refresher: what is a namespace?

_from the Linux `namespaces(7)` man page:_
> A namespace wraps a global system resource in an abstraction that makes it appear to the processes within the namespace that they have their own isolated instance of the global resource. Changes to the  global resource are visible to other processes that are members of the namespace, but are invisible to other processes.

| Namespace | Isolates |
|:----------:|----------|
| IPC | System V IPC, POSIX message queues |
| Network | Network devices, stacks, ports, etc. |
| Mount | Mount points |
| PID | Process IDs |
| User | User and group IDs |
| UTS | Hostname and NIS domain name |

_Docker has supported all of these more-or-less since the beginning, except `user`, added in v1.10 (early 2016)._

---

## User namespace remapping

- runs all containers with a different `user` namespace from the host
- UIDs, which can be `0` through `65536`, are _mapped_ to much higher values
  - i.e. `0` in a remapped container becomes e.g. `100000`, etc...

```console
$ docker run -it --rm alpine
/ # id
uid=0(root) gid=0(root) groups=0(root),1(bin),2(daemon),3(sys),4(adm),6(disk),10(wheel),11(floppy),20(dialout),26(tape),27(video)
```
```console
$ ps faxu
root      7734  0.1 13.3 337564 81112 ?        Ssl  00:41   0:09 /usr/bin/dockerd -H fd://
root      7764  0.1  3.9 293412 24180 ?        Ssl  00:41   0:15  \_ docker-containerd --config /var/run/docker/containerd/containerd.toml
root     15661  0.0  0.5   7648  3516 ?        Sl   03:01   0:00      \_ docker-containerd-shim -namespace moby -workdir /var/lib/docker/296608.296608/containerd/daemon/io.containerd.runtime.v1.linux/moby/ebf29306185c4e3ce1003a78f70bf1132
296608   15676  0.0  0.0   1584     4 pts/0    Ss+  03:01   0:00          \_ /bin/sh
```

Notice the UID is `296608` on the host, but `0` in the container!

---

## Enabling userns-remap

- _opt-in_ for compatibility reasons
- as a `dockerd` argument - `--userns-remap="someuser"`
- in the `/etc/docker/daemon.json` file (preferred):
```json
{
  "userns-remap": "testuser"
}
```
- use the value `default` to have Docker set up a `dockremap` map for you
- see [docs](https://docs.docker.com/engine/security/userns-remap/) for details!

---

## Caveats

- only one user namespace can be used - shared by all containers run by the daemon (due to kernel limitations, see also [moby/moby#28593](https://github.com/moby/moby/issues/28593))
- you can still escape the namespace with `docker run --userns=host`
  - if someone can run `docker run`, they effectively already have `root`
- using `--pid=host` or `--network=host` doesn't work properly
- using `--privileged` doesn't work
- file ownership in mounted volumes will need some forethought
- `mknod` (device creation) won't work, even as `root` within the container, even with `CAP_MKNOD` capability provided

---

## References...

- http://canihaznonprivilegedcontainers.info
- https://docs.docker.com/engine/security/userns-remap/

---

## Thank you!

### Any Questions?

