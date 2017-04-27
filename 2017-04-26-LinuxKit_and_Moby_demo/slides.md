layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# LinuxKit and Moby

_Dave Henderson, April 2017_

---

## LinuxKit - what is it?

_A toolkit for building secure, portable and lean operating systems for
containers_

- Inspired by from the experience of building Docker Editions
  - _e.g. Docker for Mac/Windows/AWS/Azure/GCP/etc..._
- Secure defaults
- Highly customizable
- Tooling intended to be _easy_
- Built with containers, for running containers

---

## LinuxKit - how do I use it?

- Simple YAML-format config file - defines what ends up in the resulting image
  - `kernel` - which kernel to use (distributed in a container image)
  - `init` - processes to start on boot (distributed in a container image, _but not run in containers_)
  - `onboot` - containers to start on boot (in order - usually system services)
  - `services` - services to run after booted
  - `trust` - lists which components (images) should be verified with Docker Content Trust prior to pulling
  - `output` - specifies what format to output (VHD, VMDK, ISO, etc...)

---

## LinuxKit - minimal example

Contains a minimal set of services: `init`, `runc` (for running containers),
`containerd` (for managing running containers), and `dhcpcd` (for getting an IP).

```yml
kernel:
  image: "linuxkit/kernel:4.9.x"
  cmdline: "console=ttyS0 console=tty0 page_poison=1"
init:
  - linuxkit/init:42fe8cb1508b3afed39eb89821906e3cc7a70551
  - linuxkit/runc:b0fb122e10dbb7e4e45115177a61a3f8d68c19a9
  - linuxkit/containerd:60e2486a74c665ba4df57e561729aec20758daed
onboot:
  - name: "dhcpcd"
    image: "linuxkit/dhcpcd:48e249ebef6a521eed886b3bce032db69fbb4afa"
    binds:
     - /var:/var
     - /tmp/etc:/etc
    capabilities:
     - CAP_NET_ADMIN
     - CAP_NET_BIND_SERVICE
     - CAP_NET_RAW
    net: host
    command: [ "/sbin/dhcpcd", "--nobackground", "-f", "/dhcpcd.conf", "-1" ]
trust:
  image:
    - "linuxkit/kernel"
outputs:
  - format: "kernel+initrd"
```

---

## `moby` - what is it?

- a tool for building LinuxKit images
- Not to be confused with "The Moby Project"
  - or "moby" which is the (temporary) name of the Docker Editions Linux distro
  - _or the whale..._
- There is also a tool named `linuxkit` used to run VMs produced by `moby build`
  - for running VM images, and pushing (to remote image stores)

---

## Demo - let's build and run a minimal image!

_Using the `minimal.yml` example in https://github.com/linuxkit/linuxkit/tree/master/examples_

1. First, we build:
  ```console
  $ moby build minimal.yml
  ...
  Create outputs:
    minimal-bzImage minimal-initrd.img minimal-cmdline
  ```
2. Now we can run:
  ```console
  $ linuxkit run hyperkit ./minimal
  ...
  / # pstree
  init-+-containerd
       |-sh---pstree
       `-sh
  ```
3. `halt` to exit the VM...

---

## The future!

- all Docker Editions will be rebased to LinuxKit-built images
- other projects (such as RancherOS) have committed to using LinuxKit as their
  new base

---

## Thank you!
