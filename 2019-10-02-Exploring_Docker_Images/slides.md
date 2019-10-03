layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-position: right 12px bottom 40px
background-size: 64px

---
class: center, middle

# ⛏ Exploring Docker Images: _Layer by layer_ 🕳

_Dave Henderson, October 2019_

???

In this talk, we'll take a deep-dive into an important aspect of Docker: images.

We'll talk about what images actually are, the different ways to make and share images, and where and how they're stored.

Whether you're brand new to Docker, or have been using it since the beginning, hopefully you'll learn something tonight!

---

## First, some basics...

```console
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
Digest: sha256:b8ba256769a0ac28dd126d584e0a2011cd2877f3f76e093a7ae560f2a5301c00
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
[...]
```

What's actually happening here? Let's break it down!

---

## First, some basics...

```console
$ docker run hello-world
```

- 😀: _"Hey Docker! I want to create and start a new container based on an image named `hello-world`!"_

---

## First, some basics...

```console
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
```

- 🤖: _I'd love to, but I don't have an image like that!_

Notice:
- the `:latest` tag is automatically appended (because another one wasn't given)

---

## First, some basics...

```console
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
```

- 🤖: _I mean, I guess I could try to download it from a registry..._

--

- ✨ _Oh and I'm assuming what you really meant was the `docker.io/library/hello-world:latest` image..._

---

## First, some basics...

```console
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
Digest: sha256:b8ba256769a0ac28dd126d584e0a2011cd2877f3f76e093a7ae560f2a5301c00
```

- 🤖: _Got it! And here's the SHA256 hash in case you want to verify_

---

## First, some basics...

```console
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
Digest: sha256:b8ba256769a0ac28dd126d584e0a2011cd2877f3f76e093a7ae560f2a5301c00
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
[...]
```

- 🤖: _All done!_

---

## Not so fast!

_I have so many questions..._

- Where did the image actually come from?
- And where did it get stored?
- And what _is_ an image anyway?

---

## Docker Registries

`docker pull` will first download a given image from a _registry_, and `docker run` will also try when the image isn't already present.

- Many hosted registries, like DockerHub (the default), or Quay.io, GCR, ECR, Artifactory, etc...
- Can also self-host:
    ```console
    $ docker run -p 5000:5000 registry
    ```

---

## Image storage

- Stored (by default) in the Docker root dir (`/var/lib/docker`)
  - see `/var/lib/docker/image/overlay2/`

- Stored as independent layers for space efficiency
  - two images that reference the same layer won't take twice the space

--

- ⚠️ _There isn't a single "image file" that you can access and share!_

---

## What's in an image?

_Generally..._

- Layers (files)
  - base OS files (packages, libraries, etc...)
  - added files (your code/binaries, etc...)
- Metadata (environment variables, commands, etc)

_But before we dive deeper..._

---

## Let's go on a tangent!

### How else can we get images?

- By building them from a `Dockerfile`:

```docker
FROM alpine:3.10

RUN apk add --no-cache figlet

COPY --from=hairyhenderson/gomplate:slim /gomplate /bin/gomplate

ENV CITY=Ottawa

CMD gomplate -i 'Hello, {{ getenv `CITY` }}' | figlet
```

```console
$ docker build -t myimage .
...
```

---

## Let's go on a tangent!

### How else can we get images?

- By committing from a container:

```console
$ docker run --name mycontainer alpine:3.10 touch /hello_docker
$ docker commit mycontainer myimage
```

---

## Let's go on a tangent!

### How else can we get images?

- By loading them:

```console
$ docker load -i image.tar
```

---

## Let's go on a tangent!

### How else can we get images?

- By importing them:

```console
$ docker import image.tar myimage
```

---

## Let's go on another tangent!

### Where did those image tar files come from!?

- _I lied earlier..._ 😅
- You _can_ make a single "image file" that you can access and share!

- Use `docker save`:
  ```console
  $ docker save -o image.tar myimage
  ```
- Or `docker export`:
  ```console
  $ docker export -o image.tar mycontainer
  ```

---

## `docker save` vs. `export`

- `save`:
  - saves images as full-fidelity archives
  - produces a difficult-to-read archive very similar to the on-disk storage format
  - can be used to share images without registries, with no loss of data or metadata

--

- `export`:
  - dumps a container's filesystem
  - produces a "regular" archive that can easily be read & explored
  - loses all metadata and potentially some data
  - more space-efficient, though!

---

## `docker save` vs. `export`

### When to use `save`?

- air-gapped or highly restricted networks where a registry can't be run

--

- temporary/local image caching (during CI, etc)

---

## `docker save` vs. `export`

### When to use `export`?

- development/troubleshooting (see also `docker cp`)

--

- ...probably never 😉

---

## `docker load` vs. `import`

- `load`:
  - restores an image completely
  - can only be used with an archive produced by `save`
  - doesn't allow image name to change

--

- `import`:
  - creates a bare-bones image with _only_ a filesystem
  - can be used to create stripped-down base images (see also `FROM scratch` with `ADD`/`COPY`)

---

## Anatomy of an image archive

- list of images in the archive (`repositories`)
- manifest (`manifest.json`)
- image config (`<image ID>.json`)
- layer directories (`<layer ID>/`)
  - `VERSION` file
  - layer config (`json`)
  - layer archive (`layer.tar`)

---

## Image Manifests

- Enables multi-platform images
- Links to specific imaged by SHA256 digests
- Use `docker manifest inspect` to see manifests

---

## A short digression: OCI Images

- Vendor-neutral image format
- Supported natively by many container runtimes (runC, runV, ClearContainers, gVisor, Kata-Containers, etc...)
- Can be created with many tools, especially `img` from Jess Frazelle

---

## Demo time! 🎭

???

- loading an image from a tar archive
- creating an image from a tar archive (with `import`)
- exploring the contents of a `save`d archive
  - note relationships to parents
  - explore config
    - note empty layers
    - note layer SHA sums
    - note config (env, cmd)
- create the same image with `img`:
    ```console
    $ docker run --privileged -it --rm -v $(pwd)/:/tmp --entrypoint /bin/sh r.j3ss.co/img
    / $ img build -t myimage /tmp
    / $ img save myimage -o out.tar
    / $ exit
    ```
  - explore `out.tar`, notice `oci-layout`

---

## Thank you! 🙇‍

### Time For Questions!
