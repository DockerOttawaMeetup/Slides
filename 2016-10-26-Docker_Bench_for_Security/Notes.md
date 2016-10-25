layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# Docker Bench for Security

_Dave Henderson, October 2016_
---

## First, some background...

---

### Isn't Docker Secure?

- Well sure (more or less), but it never runs in a vacuum!

--
- Needs a host OS to run on

--
- Pointless without containers

???
Docker has lots of security features, and on its own is _pretty_ secure
out of the box, but it isn't an OS, and is pointless without containers.

---

### But I know nothing about security!

--

- Neither do I 😜
--

- Thankfully, other people do!
--

- Specifically, the _CIS_

![CIS Logo](https://www.cisecurity.org/images/CISLogoWeb.png)

---

### What does the CIS do?

- Independent security reviews - i.e. they'll (probably) be more _honest_
- Security [Benchmarks](https://benchmarks.cisecurity.org)
  - especially the [CIS Docker Benchmark](https://benchmarks.cisecurity.org/tools2/docker/CIS_Docker_1.12.0_Benchmark_v1.0.0.pdf)
  - created by CIS together with Docker and

---

### Let's look at the benchmark...

- Set of _guidelines_, some "Scored", some "Not Scored" (~optional, or hard to
  test automatically)
- Separate focus on:
  - Host OS configuration
  - Docker Engine configuration
  - Container configuration
  - Ops considerations
--

- Let's dig in to the [CIS Docker Benchmark](https://benchmarks.cisecurity.org/tools2/docker/CIS_Docker_1.12.0_Benchmark_v1.0.0.pdf)

???
When going through the PDF, probably only useful to look at the ToC and click
through to a few interesting spots.
---

## Making things easier...

- Going through and validating each of these guidelines would take _hours_...
--

- So Docker's security team has automated it for us!

---

### Docker Bench for Security

> _a script that checks for dozens of common best-practices around deploying Docker containers in production_

- Just run it and the CIS Docker Benchmark best
  practices are checked
- _Fairly_ easy to run...
  ```console
  $ docker run -it --net host --pid host --cap-add audit_control \
    -v /var/lib:/var/lib \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /usr/lib/systemd:/usr/lib/systemd \
    -v /etc:/etc --label docker_bench_security \
    docker/docker-bench-security
  ```
  - or you can run their script on your host:
    ```console
    $ git clone https://github.com/docker/docker-bench-security.git
    $ cd docker-bench-security
    $ sh docker-bench-security.sh
    ```

???
Even though there's a CIS Benchmark for Docker 1.12, Bench hasn't been updated
yet. It's still targeting the Benchmark from Docker 1.10.
---

class: middle, center

# Demo time!

???
deploy with docker-machine
```
docker-machine create -d amazonec2 --amazonec2-zone=b benchmachine
eval $(docker-machine env benchmachine)
```
start a sample workload (docker/example-voting-app)
```
docker-compose -f example-voting-app/docker-compose.yml up -d
```
run docker-bench containerized
```
docker run -it --net host --pid host --cap-add audit_control \
  -v /var/lib:/var/lib \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /usr/lib/systemd:/usr/lib/systemd \
  -v /etc:/etc --label docker_bench_security \
  docker/docker-bench-security
```
inspect results...

now run the local script - see how this has a bit more?
```
docker-machine ssh benchmachine git clone https://github.com/docker/docker-bench-security.git
docker-machine ssh benchmachine sudo apt-get update
docker-machine ssh benchmachine sudo apt-get install auditd
docker-machine ssh benchmachine "cd docker-bench-security; sudo ./docker-bench-security.sh"
```
---

## Where next?

- Fix all the things!
--

- You may need a custom VM image for your distro of choice
--

- Some issues are fairly simple to resolve
  - 4.1  - Create a [non-`root`] user for the container
  - 5.10 - Limit memory usage for container
  - 5.11 - Set container CPU priority appropriately
  - 5.13 - Bind incoming container traffic to a specific host interface
--

- Some issues will take more effort...
  - 4.5  - Enable Content trust for Docker
  - 5.1/5.2 - SELinux/AppArmor verification
--

- Some issues don't _need_ to be fixed for your usage!
  - Maybe you _want_ the host network exposed to some of your containers (5.13)
  - Maybe you _need_ to access a host device in the container (5.17)
--

- Remember to re-run periodically, ideally in a CI/CD workflow
  - some related tools that may help:
    - Actuary - https://github.com/diogomonica/actuary
    - Docker Bench Test - https://github.com/gaia-adm/docker-bench-test
--

- Use your judgement, and remember that these guidelines aren't a silver security
  bullet.

---

## Thank you!
