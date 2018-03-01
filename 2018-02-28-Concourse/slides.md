layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# 🛫 Playing with Concourse CI 🛬

_Dylan Arbour, February 2018_

---

## What is Concourse CI? ✈️

- Started at Pivotal (RabbitMQ, CloudFoundry)
- Open Source
- Written in Golang
- Every job/task runs in a container
- Pipelines are first class citizens
- Pipelines are defined as code (yaml)
- Can be controlled by CLI (fly)

https://github.com/concourse/concourse

---

## Pipelines 🚰

- Similiar to Circle CI's 2.0 Workflows
- Many gates, forks, rules
- Composed of 3 parts

--

```
                  +------------------------+
                  |                        |
                  |                        v
       X          |   +----------------------------------------+
      / \         |   |                                        |
     /   \        |   |  Job                                   |
    /     \       |   |                                        |
   /       \      |   |  +---------------------------------+   |
  /         \     |   |  |                                 |   |
 /           \    |   |  |                                 |   |
|   Resource  |---+   |  |                                 |   |
 \           /        |  |             Task                |   |
  \         /         |  |                                 |   |
   \       /          |  |                                 |   |
    \     /           |  |                                 |   |
     \   /            |  |                                 |   |
      \ /             |  +---------------------------------+   |
       V              +----------------------------------------+
```
---

## Resources 🚗

- Examples: Git repository, Docker Image
- Lots of custom options (Slack notifications, other pipelines)
- Can be pulled or pushed
- Can watch for changes

--

```yaml
resources:
- name: hello-world
  type: docker-image
  source:
    repository: hello-world
    tag: latest
```

---

## Jobs 🔧

- Can have many tasks
- Can run sync or async
- Can be triggered by a resource change

--

```yaml
jobs:
- name: speak
  plan:
  - get: my-resource
    trigger: true
  - task: say-hello
    file: ./say-hello.yaml
  - task: say-goodbye
    file: ./say-goodbye.yaml
```

---

## Tasks 🔩

- The worker
- Runs in a container
- Can push and pull a resource
- Has contraints like `passed`, `failed`

--

```yaml
- task: say-hello
  config:
    platform: linux
    image_resource:
      type: docker-image
      source: {repository: ubuntu}
    run:
      path: echo
      args: ["Hello, world!"]
```

---

## Demo 👩‍💻

---

## More Examples 📚

- http://concourse.ci/flight-school.html
- https://github.com/pivotalservices/concourse-pipeline-samples
- https://github.com/starkandwayne/concourse-tutorial

---


## Thank you! 🙇‍♂️

### Questions?
