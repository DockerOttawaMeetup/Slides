layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# Creating Multi-Arch Docker Images

_Dave Henderson, February 2018_

---

## What is Multi-Arch?

- Being able to run the _"same"_ image on any of the OSes and architectures
  that Docker supports:
  - `windows`, `linux`
  - `arm64`, `arm`, `ppc64le`, `amd64`, `i386`, `s390x`...
- DockerHub updated in September 2017 to support multi-arch, and many
  official DockerHub images were made multi-arch

---

## The old and hacky way 😩

- uses a prefix to indicate architecture for the "other" arches:
    ```console
    $ docker run arm32v7/ubuntu uname -sm
    Linux armv7l
    ```
- this is bad because it requires users to remember inconsistent arch names
    ```console
    $ docker run arm32v7/alpine uname -a
    ...
    docker: Error response from daemon: pull access denied for arm32v7/alpine,  repository does not exist or may require 'docker login'.
    ```
    - _oh right..._
    ```console
    $ docker run armhf/alpine uname -a
    Linux armv7l
    ```

---

## The new and awesome way! ✨

- By _"same"_ I mean the same image _name_
  - The command `docker run ubuntu` should work the same on Raspberry Pi as it
    does on an x86-64 laptop, or a System Z mainframe
    ```console
    $ docker run alpine uname -a
    Linux 9d0c6f913736 4.9.59-v7+ #1047 SMP Sun Oct 29 12:19:23 GMT 2017 armv7l Linux
    ```

---

## How does this work? 🔎

- On a `docker push`, a _manifest_ is created on the registry to act as a
  pointer to image metadata and layers
- Support for _manifest lists_ (a.k.a. _"fat manifests"_) added in Docker
  1.10 (Jan 2012), allowing to group many manifests in a list
  - manifest entries contain a `platform` section:
      ```json
      "platform": {
          "architecture": "ppc64le",
          "os": "linux"
      }
      ```
- On a `docker pull`, either a manifest or manifest list is returned, and
  Docker decides which image to pull based
  on the current platform

---

## New `manifest` command 🚀

- Docker 18.02 now contains an _experimental_ `manifest` command:
  ```console
  $ docker manifest --help
  Usage:	docker manifest COMMAND
  ...
  Commands:
    annotate    Add additional information to a local image manifest
    create      Create a local manifest list for annotating and pushing to a registry
    inspect     Display an image manifest, or manifest list
    push        Push a manifest list to a repository
  ```
- Replaces an earlier stand-alone [`manifest-tool`](https://github.com/estesp/manifest-tool)
  by Phil Estes

---

## Creating manifest lists 📦

- First, let's build a bunch of cross-platform images
  - _see `code/build_linux.sh`_
- Because the `docker manifest` commands operate on image manifests pulled from a registry, we need to push these:
  - `docker push hairyhenderson/helloarch`
- Now create the manifest list with `docker manifest create MANIFEST_LIST MANIFEST [MANIFEST...]`
  - _see `code/manifest_list.sh`_
- And make sure it was created with `docker manifest inspect MANIFEST_LIST`

---

## Patching things up 🛠

- We built all images on the same system, which means the image OS and
  Archictecture metadata will be set to `linux`/`amd64`.
- We need to run `docker manifest annotate` for each other architecture
  to apply the correct metadata 
  - _(see `manifest_linux.sh`)_
- Finally, we push the manifest list
  - `docker manifest push hairyhenderson/helloarch:latest`
- Now we can run the same image on different architectures!

---

## How about Multi-Arch _and_ Multi-OS?

- Windows image needs to be built on Windows
  - _(see `code/build_win.ps1`)_
- We must push the `windows-amd64` image, and then re-create and push
  the manifest list
  - _(see `code/manifest_win.ps1`)_

---

## The payoff! 💰

```console
linux1@zvm:~> docker run --rm hairyhenderson/helloarch
Hello, linux/s390x!
```

```console
hairyhenderson@rpi3:~$ docker run --rm hairyhenderson/helloarch
Hello, linux/arm!
```

```console
PS C:\> docker run --rm hairyhenderson/helloarch
Hello, windows/amd64!
```

---

## Digging Deeper 🕳

- [`docker manifest` docs](https://docs.docker.com/edge/engine/reference/commandline/manifest/)
- [Docker Image Manifest v2.2 spec](https://github.com/docker/distribution/blob/master/docs/spec/manifest-v2-2.md), also see [OCI Image Index Spec](https://github.com/opencontainers/image-spec/blob/v1.0.1/image-index.md)
- ["From Arm to Z" talk at DockerCon 2017](https://www.youtube.com/watch?v=nrBYUw1Pz5I) _([slides](https://www.slideshare.net/Docker/from-arm-to-z-building-shipping-and-running-a-multiplatform-docker-swarm-christy-perez-and-christopher-jones-ibm))_
  - _from Christy Perez, author of the `docker manifest` command, and Christopher Jones_
- [PR #138 on docker/cli](https://github.com/docker/cli/pull/138)
  - _the PR that added the feature_
- [Docker Blog: "Docker Official Images Are Now Multi-Platform"](https://blog.docker.com/2017/09/docker-official-images-now-multi-platform/)
- [Docker BLog: "Multi-Arch All The Things](https://blog.docker.com/2017/11/multi-arch-all-the-things/)

---

## Thank you!

### Any Questions?

