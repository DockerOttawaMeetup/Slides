layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---

class: middle, center
# A Docker Clustering Example Using RabbitMQ
_Etienne Dufresne, August 2016_

---

# GitHub Repo

https://github.com/EtienneDufresne/rabbitmq-docker-cluster

---

## What is RabbitMQ?

--

###Messaging Service

--

Supports Work Queues

![alt text](https://www.rabbitmq.com/img/tutorials/python-two.png)

--

Supports Topic

![alt text](https://www.rabbitmq.com/img/tutorials/python-five.png)

--

Can be Clustered for speed and high availability

---

## Requirements for the RabbitMQ Cluster

--

- All cluster members should run in its own Docker container

--

- A Docker host may run a single or several RabbitMQ services

--

- The first RabbitMQ service to come up becomes master

--

- Subsequent RabbitMQ services become slaves

--

- If several RabbitMQ services come up at the same time the cluster should be formed without issues

--

- RabbitMQ requests should be load balance across RabbitMQ services

--

- If a RabbitMQ service crashes, we should not lose any data

--

- If the master RabbitMQ service goes down, a slave should be promoted as the new master

--

- If a starting or restart RabbitMQ service comes up it should be able to join the cluster

--

- We need to be able to monitor the health of every RabbitMQ service

---

## Problems to Solve

--

How are the Docker containers running RabbitMQ are going communicate?

--

By using a user defined Docker network

--

How can I ensure I'm not going to lose data if I'm running in a container?

--

By using a Docker volume mounted to the Docker host's disk

--

If I am a new RabbitMQ service how can I tell if other services are already clustered?

--

By using a centralized repository of services that can be used to discover services at runtime

--

How can I tell if my cluster is healthy?

--

By registering health checks in the service repository

---


## User Defined Network

docker-compose version 2 allows the definition of user defined networks in a top level section
![alt text](https://docs.docker.com/engine/userguide/networking/images/bridge_network.png)

--

```YML
version: '2'
networks:
  rabbitmq_network:
    driver: bridge
    ipam:
      driver: default
      config:
      - subnet: 172.20.0.0/24
        gateway: 172.20.0.1
```

--

This network provides a DNS to resolve container IPs by name

---

## Docker Volumes

docker-compose version 2 allows the definition of volumes in a top level section

--

```YML
volumes:
  rabbitmq_persistence1: {}
  rabbitmq_persistence2: {}
```

---

## Service Registry and Discovery

###Consul by HashiCorp

https://www.consul.io/

--

- Allows the registrations of services and their health checks into a centralized repository

--

- The repository can be queried by DNS or HTTP to obtain service information

--

- When queried for a service, Consul API load balances across healthy services

--

- Provides a general purpose key value store that can be used to hold leader elections

---

## Service Registry and Discovery

```YML
consul-server:
  image: gliderlabs/consul-server:latest
  restart: always
  container_name: consul
  command: -bootstrap-expect 1
  environment:
    - SERVICE_8500_NAME=consul
    - SERVICE_8600_NAME=consul
  ports:
    - "8500:8500"
    - "8600:8600"
  networks:
    rabbitmq_network:
      ipv4_address: 172.20.0.10
```

--

Joins the rabbitmq_network and is given a static ip

--

Ports 8500 and 8600 are exposed

---

## Consul's DNS API

Runs on port 8600 but it needs to be made available to other containers on port 53

--

Use dnsmask to forward requests

--

```YML
dnsmask:
  image: andyshinn/dnsmasq:latest
  depends_on:
    - consul
  restart: always
  container_name: dnsmask
  environment:
    - SERVICE_53_NAME=dnsmask
    - SERVICE_TAGS=dnsmask
  ports:
    - 53:53/tcp
    - 53:53/udp
  cap_add:
    - NET_ADMIN
  command: -S /consul/172.20.0.10#8600 --log-facility=-
  networks:
    rabbitmq_network:
      ipv4_address: 172.20.0.11
```
--

Joins the rabbitmq_network and is given a static ip

--

Other containers must set their DNS to 172.20.0.11

---

## Automatic Service Registration

We don't want to have to manually register and deregister services with Consul as they come up and down.

--

```YML
registrator:
  image: gliderlabs/registrator:latest
  depends_on:
    - consul
  restart: always
  container_name: registrator
  volumes:
    - /var/run/docker.sock:/tmp/docker.sock
  command: consul://consul:8500
  networks:
    - rabbitmq_network
```

--

Registrator also supports automatic health check registration
---

### Let's docker-compose our way out!

Start consul, dnsmask and registrator

```shell
docker-compose up -d consul dnsmask registrator
```

--

Let's ensure the Docker network DNS is working

```shell
docker run -it --network rabbitmqdockercluster_rabbitmq_network --dns 172.20.0.11 debian ping -c 1 consul
```

--

Let's ensure the Consul DNS API is working on port 53 (via dnsmask)

```shell
docker run -it --network rabbitmqdockercluster_rabbitmq_network --dns 172.20.0.11 debian ping -c 1 consul.service.consul
```

--

Let's ensure the Consul HTTP API is working on port 8500

http://localhost:8500/ui

---

### Diving into the RabbitMQ Dockerfile

Extends the official RabbitMQ management Docker image and overrides the entrypoint

--

```YML
FROM rabbitmq:3.6.5-management

ENV RABBITMQ_DEFAULT_USER=guest
ENV RABBITMQ_DEFAULT_PASS=guest
ENV RABBITMQ_ERLANG_COOKIE=supersecretcookie
ENV SERVICE_15672_IGNORE=true
ENV SERVICE_5672_IGNORE=true

RUN apt-get update && apt-get install -y jq curl dnsutils

COPY docker-entrypoint.sh /usr/local/bin/
```

---

### Diving into the entrypoint

--

### AKA let's Bash our way out

--

It's a copy of the official RabbitMQ entrypoint that performs additional tasks

--

It uses the Consul key value HTTP API to obtain a lock (semaphore)

--

If the lock can't be obtain it waits until another instance is done using it

--

It generates a list of healthy available RabbitMQ services using dig and the Consul DNS API

```shell
dig rabbitmq.service.consul
```

--

It writes a RabbitMQ config file with the available RabbitMQ services available for clustering

--

It registers the RabbitMQ service with Consul along with a health check

---

### Let's docker-compose our way out

--

###... AGAIN!

--

Build the container

```shell
docker build -t rabbitmq-docker-cluster .
```

--

Start two RabbitMQ services to test our leader election logic

```shell
docker-compose up -d rabbitmq1 rabbitmq2
```

Let's look at the output from the two containers

```shell
docker logs -f rabbitmq1
```

```shell
docker logs -f rabbitmq2
```

---

### Let's make sure it worked

Let's ensure that the services are registered and healthy in Consul

http://localhost:8500/ui

--

Let's ensure that the RabbitMQ management interface is also seeing the cluster

http://localhost:15672/#/

--

Let's run a producer / consumer test

```shell
npm install
```

```shell
node consumer.js
```

```shell
node producer.js
```

---

### Let's test some hight availability scenarios

With the consumer stopped, let's add a few messages to the queue

```shell
node producer.js
```

--

Stop the master RabbitMQ service and make sure the queue is still available

```shell
docker-compose stop rabbitmq1
```

http://localhost:15673/#/queues/%2F/my-queue

--

Stop the slave RabbitMQ service and check that the console is down

```shell
docker-compose stop rabbitmq2
```

---

### Let's test some hight availability scenarios

Start the RabbitMQ slave service and notice that we did not lose any data

```shell
docker-compose up -d rabbitmq2
```

--

Start the RabbitMQ master service and notice that the queue is in sync

```shell
docker-compose up -d rabbitmq1
```

http://localhost:15672/#/queues/%2F/my-queue

---

### Did we meet the requirements?

- ✓ All cluster members should run in its own Docker container
- ✓ A Docker host may run a single or several RabbitMQ services
- ✓ The first RabbitMQ service to come up becomes master
- ✓ Subsequent RabbitMQ services become slaves
- ✓ If several RabbitMQ services come up at the same time the cluster should be formed without issues
- ✓ RabbitMQ requests should be load balance across RabbitMQ services
- ✓ If a RabbitMQ service crashes, we should not lose any data
- ✓ If the master RabbitMQ service goes down, a slave should be promoted as the new master
- ✓ If a starting or restart RabbitMQ service comes up it should be able to join the cluster

--

Yes but this is on a single Docker host, what if it goes down?

---

### Making it work on multiple Docker hosts

RabbitMQ wants resolvable hostnames but we do not want to have to give static IPs and hostnames to our RabbitMQ services as we want the cluster to be able to scale dynamically

--

The current docker-compose file can work on multiple Docker hosts instances with a few tweaks

--

Use https://github.com/gliderlabs/hostlocal and host networking

--

This approach only supports a single RabbitMQ service per host as only one RabbitMQ container can use the host's IP and hostname

---

### The real multi host solution

Use Docker Swarm

--

![alt text](https://docs.docker.com/engine/userguide/networking/images/overlay-network-final.png)

- Swarm mode allows the definition of overlay networks which span across multiple hosts
- The overlay network's DNS can resolve containers running on any Docker host

--

- The new Swarm mode implements its own service registry (orchestration)
- It introduces the `docker service` command and `bundles` (beta)

https://blog.docker.com/2016/06/docker-1-12-built-in-orchestration/

---
class: middle, center
# Thank You!
