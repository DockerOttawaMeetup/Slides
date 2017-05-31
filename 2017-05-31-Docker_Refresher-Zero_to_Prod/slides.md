layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# Docker Refresher!

## _From Zero to Production with Docker..._

_Dave Henderson, May 2017_

---

## What is Docker?

### The high-level pitch

> _Docker is an open platform for developing, shipping, and running applications._
> _Docker enables you to separate your applications from your infrastructure so you can deliver software quickly._
> _With Docker, you can manage your infrastructure in the same ways you manage your applications._
> _By taking advantage of Docker’s methodologies for shipping, testing, and deploying code quickly, you can significantly reduce the delay between writing code and running it in production._

_(https://docs.docker.com/engine/docker-overview/)_

???

1. an open-source tool useful to dev and ops (and everyone else involved in
  shipping software)
2. extends the concept of _"immutable infrastructure"_ to package software once
  and run anywhere, while also allowing you to configure network and other
  resources specifically for that software
3. Docker solves the "works on my machine" syndrome by packaging all dependencies
  with the application

---

## What is Docker?

### Containers

- namespaces: isolated execution scope for processes
  - things like process, network, filesystem isolation
- control groups (cgroups): resource access limits
  - allows setting limits on things like RAM/CPU usage
  - limits how the process can interact with the kernel (system calls)

---

## What is Docker?

### Container Images

- a container's state on disk
- layered by use of "union" or "copy-on-write" filesystems
  - allows far more efficient storage and distribution, only need to store a
    layer once, even if shared by many Images
  - represents the state of the container at a point in time
- can be defined in `Dockerfile`s or by taking snapshots of running containers

---

## What is Docker?

### The `Dockerfile`

- a build file for container images
- defines each layer in a separate line
- uses instructions like `RUN`, `COPY`, `ENV` to build container filesystem and
  set run-time metadata

```Dockerfile
FROM debian:jessie

RUN apt-get update
RUN apt-get install figlet

ENV THING=World

CMD figlet Hello, $THING
```
_A naïve Dockerfile_

(let's demo this now!)

---

## What is Docker?

### Products

- Docker - the CLI and engine
  - runs and manages containers on a host (physical or VM)
- Docker Compose - a tool for defining and running multi-container applications
- Docker Swarm - clustering for Docker containers, embedded by default in Docker
  - runs and manages containers across multiple hosts
- Docker Registry - distribution service for container images
  - open-source for self-managed registries
  - DockerHub and Docker Store are public/hosted registries
  - 3rd-party registries like Quay.io and Amazon ECR

---

## Dockerizing an app

### Introducing `restdemo`

- simple RESTful app built with [goa](https://goa.design/)
- pretty much no UI, intended only for use with `curl` or `restdemo-cli`
- we want to deploy this as a highly-available service in the cloud, with
  a health check

---

### Creating a docker image

- the [`golang:1.8`](https://hub.docker.com/_/golang/) image in DockerHub gives us all we need to build a Go app
  ```dockerfile
  FROM golang:1.8

  RUN mkdir -p /go/src/github.com/hairyhenderson/restdemo
  WORKDIR /go/src/github.com/hairyhenderson/restdemo
  COPY . /go/src/github.com/hairyhenderson/restdemo

  RUN go build -ldflags "-w -s"

  CMD [ "/go/src/github.com/hairyhenderson/restdemo/restdemo" ]
  ```
- now we can build it:
  ```console
  $ docker build -t hairyhenderson/restdemo:v1 .
  ```
- and run it:
  ```console
  $ docker run -p 8000:8000 -it hairyhenderson/restdemo:v1
  ...
  ```

---

### Publish the image to DockerHub

- let's make the image available (public or private)
  ```console
  $ docker push hairyhenderson/restdemo:v1
  ```
- now we can see it in [DockerHub](https://hub.docker.com/r/hairyhenderson/restdemo/tags)

---

### Refactor...

- the image in DockerHub is ~270MBs - that's a bit large for a 6MB binary!
- what's worse, the size in [DockerHub](https://hub.docker.com/r/hairyhenderson/restdemo/tags) doesn't include the base `golang:1.8` image:
```console
$ docker image ls hairyhenderson/restdemo:v1     
REPOSITORY                TAG                 IMAGE ID            CREATED             SIZE
hairyhenderson/restdemo   v1                  f3912df20841        2 minutes ago       713MB
```
- Let's switch to an [Alpine Linux](https://alpinelinux.org)-based image:
```dockerfile
FROM golang:1.8-alpine
```
- looking better:
```console
$ docker image ls hairyhenderson/restdemo:v2
REPOSITORY                TAG                 IMAGE ID            CREATED             SIZE
hairyhenderson/restdemo   v2                  622a6cca2f4a        17 seconds ago      271MB
```
- it's still ~80MBs in DockerHub, but that's OK for now

---

### Let's run this thing in the cloud!

- if we use https://play-with-docker.com, we don't have to provision a "real" server (yet)!
```console
$ docker run -d -p 8000:8000 --name restdemo hairyhenderson/restdemo:v2
$ docker logs -f restdemo
...
```
- we can use `curl` now to get the `/health` endpoint
```console
$ curl http://pwd10-0-28-3-8000.host1.labs.play-with-docker.com/health
{"hostname":"09f330f871f6","version":"v2"}
```
- 🎉

---

### So what about that health check?

- we can use the `HEALTHCHECK` instruction in the `Dockerfile`
- any command can be used to check service health, as long as it exits with a `0` (healthy) or `1` (unhealthy). We'll use `curl` for now:
```dockerfile
HEALTHCHECK --interval=5s --timeout=5s \
  CMD curl --fail http://localhost:8000/health || exit 1
```
- note the interval/timeout overrides (these default to 30s)
- now, we see a health status from `docker ps`:
```console
$ docker run -p 8000:8000 -d hairyhenderson/restdemo:v3
170b67155deb7f6463fef924df43a91004766e080187a73d24a9d3703cca0a29
$ docker ps --format '{{.Image}} - {{.Status}}'
hairyhenderson/restdemo:v3 - Up 1 second (health: starting)
# wait a few seconds...
$ docker ps --format '{{.Image}} - {{.Status}}'
hairyhenderson/restdemo:v3 - Up About a minute (healthy)
```
  - 💡Pro Tip: use the `--format` flag to make the output easier to read!⭐

---

### Alright, time for some HA

- to deploy Docker containers in a self-healing HA cluster, we can use Docker Swarm
  - setting up Swarm is simple - just run `docker swarm init` to start (we can even do this in play-with-docker)
  - this requires having Docker already installed, and adds management 
  overhead
  - using one of the [Docker Editions for Cloud](https://store.docker.com/search?offering=community&platform=cloud&q=&type=edition),
  like Docker for AWS or Docker for Azure, we can set up a scalable multi-node
  Swarm with just a few clicks
- To save time, I've set up a Docker for AWS cluster in advance

---

### Deploy to a Swarm

- Swarm abstracts the concept of a container into "services"
- we use `docker service create` to create new services
```console
$ docker service create --detach --name restdemo -p 8000:8000 \
  --env VER=v3 \
  --update-delay=15s --update-parallelism=1 \
  hairyhenderson/restdemo:v3
```
- we can now access this service (first find the ELB's DNS entry from the [AWS console](https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks?filter=active&tab=outputs)):
```console
$ curl http://Docker-ExternalLoa-XXX.us-east-1.elb.amazonaws.com:8000/health
```

---

### Scale it up!

- by default, services are started in "replicated" mode (vs. "global"), with a single replica, but we can scale with:
```console
$ docker service scale restdemo=3
```
- when we list services (`docker service ls`) we can see that not all replicas are immediately available - this is because the healthcheck takes some time to execute
- `docker service ps` shows the running replicas:
```console
$ docker service ps restdemo
ID            NAME        IMAGE                        NODE                            DESIRED STATE       CURRENT STATE           ERROR               PORTS
prij602t9dn0  restdemo.1  hairyhenderson/restdemo:v3   ip-172-31-11-224.ec2.internal   Running             Running 4 minutes ago                       
hvkbyecmvd4w  restdemo.2  hairyhenderson/restdemo:v3   ip-172-31-27-140.ec2.internal   Running             Running 3 minutes ago                       
ze68ce5z0dtl  restdemo.3  hairyhenderson/restdemo:v3   ip-172-31-47-145.ec2.internal   Running             Running 3 minutes ago  
```
- Swarm will load-balance requests among service replicas (individual containers), so each `curl` command is potentially routed to a different container:
```console
$ curl http://Docker-ExternalLoa-XXX.us-east-1.elb.amazonaws.com:8000/health
{"hostname":"4670864b148b","version":"v3"}
$ curl http://Docker-ExternalLoa-XXX.us-east-1.elb.amazonaws.com:8000/health
{"hostname":"742e2a0f9f4d","version":"v3"}
$ curl http://Docker-ExternalLoa-XXX.us-east-1.elb.amazonaws.com:8000/health
{"hostname":"72e08a762c79","version":"v3"}
```

---

### Tangent: _I still don't like how big the image is..._

- the `hairyhenderson/restdemo:v3` image _is_ smaller than `v1`, but I think we can do better...
- much of the space is taken up by the Go compiler and other tools we don't need at runtime
- let's use the new _multi-stage build_ feature to split the build image from the runtime image!
- while we're at it, let's use the Goa-generated commandline client instead of `curl` for the health-check

---

```dockerfile
FROM golang:1.8-alpine AS build

RUN mkdir -p /go/src/github.com/hairyhenderson/restdemo
WORKDIR /go/src/github.com/hairyhenderson/restdemo
COPY . /go/src/github.com/hairyhenderson/restdemo

RUN go build -ldflags "-w -s"
RUN go build -ldflags "-w -s" -o cli github.com/hairyhenderson/restdemo/tool/restdemo-cli/

FROM alpine:3.6

COPY --from=build /go/src/github.com/hairyhenderson/restdemo/restdemo /app/restdemo
COPY --from=build /go/src/github.com/hairyhenderson/restdemo/cli /app/cli
COPY --from=build /go/src/github.com/hairyhenderson/restdemo/public /app/public

WORKDIR /app
HEALTHCHECK --interval=5s --timeout=5s CMD ./cli show health || exit 1

CMD [ "./restdemo" ]
```

- now the size is _much_ smaller!
```console
$  docker images hairyhenderson/restdemo:v4
REPOSITORY                TAG    IMAGE ID       CREATED        SIZE
hairyhenderson/restdemo   v4     5fb33400cbda   5 seconds ago  16.3MB
```

---

### Upgrade time!

- ok, let's build this new image as `:v4` and do a rolling upgrade!
```console
$ docker service update --image hairyhenderson/restdemo:v4 restdemo
```
- and we can use `docker service ps` to watch the versions roll!
  - again, we can use the `--format` option to present a different picture:
```console
$ docker service ps --format '{{.Name}} - {{.Image}} - {{.CurrentState}}' restdemo 
restdemo.1 - hairyhenderson/restdemo:v4 - Running 3 minutes ago
restdemo.1 - hairyhenderson/restdemo:v3 - Shutdown 3 minutes ago
restdemo.2 - hairyhenderson/restdemo:v4 - Running 3 minutes ago
restdemo.2 - hairyhenderson/restdemo:v3 - Shutdown 3 minutes ago
restdemo.3 - hairyhenderson/restdemo:v4 - Running 3 minutes ago
restdemo.3 - hairyhenderson/restdemo:v3 - Shutdown 3 minutes ago
```

---

## Thank you!
