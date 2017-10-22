layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# Prometheus on Docker for AWS: _Lessons learned..._

_Dave Henderson, July 2017_

---

## Introduction

1. What's Prometheus?
2. How Qlik is using Prometheus
2. How to deploy these in Docker for AWS
3. How to have a data throughput crisis
4. How to fix your data throughput crisis!

---

## What's Prometheus?

> _Prometheus is an open-source systems monitoring and alerting toolkit originally built at SoundCloud. Since its inception in 2012, many companies and organizations have adopted Prometheus, and the project has a very active developer and user community. It is now a standalone open source project and maintained independently of any company._
>

### Features

- a multi-dimensional data model (time series identified by metric name and key/value pairs)
- a flexible query language to leverage this dimensionality
- no reliance on distributed storage; single server nodes are autonomous
- time series collection happens via a pull model over HTTP
- pushing time series is supported via an intermediary gateway
- targets are discovered via service discovery or static configuration
- multiple modes of graphing and dashboarding support

---

## Prometheus Architecture

.center[<img src="https://prometheus.io/assets/architecture.svg" style="height: 360px"/>]

Prometheus scrapes metrics from instrumented jobs, storing all scraped samples locally and running rules over this data to either record new time series from existing data or generate alerts.
Grafana or other API consumers can be used to visualize the collected data.

---

## Scaling and Federating Prometheus

- A single Prometheus instance can typically handle millions of time series - enough to scrape 1000 servers with 1000 time series each, every 10 seconds
- Prometheus is stand-alone, so you can't scale horizontally in a typical cluster fashion
- Instead, Prometheus scales through federation
- Other than scaling, federation is important for monitoring multiple datacenters or regions
- How Qlik federates between regions:

![](./images/diag1.svg)

---

## Deployment considerations - data persistence

- Prometheus is a stateful service - its data must be stored somewhere if we want it to persist across restarts
- We have 4 Prometheus instances: 3 regional "scrapers", and 1 global "federator"
  - the scraper data is not _too_ important
    - 24hr retention (`-storage.local.retention=24h0m0s`)
    - federation happens every 60s so it's OK to lose a sample or two
    - _but_ alerts are fired from these instances, so we should have enough samples to monitor accurately
  - the federator data is _pretty_ important
    - 15 day retention (default)
    - we query this instance for dashboards - good to have a reasonable amount of history available
- The way to persist data with Docker (especially in Swarm) is Volumes
  - Docker for AWS gives us choice of 3 volume drivers by default:
    - Local driver (data stored on host's disk)
    - Cloudstor driver: Shared Volumes (backed by EFS)
    - Cloudstor driver: Relocatable Volumes (backed by EBS)
- Prometheus is quite space-efficient
  - our production federator uses only 14GBs for 15 days of data on ~150k time series
- We chose EFS-backed Shared Volumes for the federator's data, and Local volumes for the target data
  - historical reasons - earlier versions of Cloudstor didn't yet support EBS
  - also EFS can be shared across Availability Zones, meaning a prometheus service instance can be scheduled to a new host in a different AZ and have access to the previous instance's data immediately

---

## Deployment considerations - prometheus configuration

- Prometheus expects a YAML-format config file, containing secrets
- We don't want to include a hard-coded config file in the Docker image
- Therefore, we generate a config file from a template at runtime
  - uses an `ENTRYPOINT` script and `gomplate`
  - reads secrets from Swarm Secrets
- The federator and target prometheus instances are different enough that it makes sense to create multiple Docker images
  - the first builds `FROM` the official `prom/prometheus` image
  - the second builds from the first, adding in a modified entrypoint script and config template

---

## Deployment considerations - Swarm

- The federator's Swarm is built around the federator, but it's not the only service
  - there's a number of supporting services also:
    - Grafana
    - MariaDB (and mysqld-exporter)
    - Redis (and redis-exporter)
    - Caddy (2 replicas)
    - node_exporter (global)
    - cAdvisor (global)
    - Alertmanager
    - and another Prometheus to monitor all of those (a.k.a `prometheus-local`)
  - we need to be realistic about resource allocation
- The target Swarm doesn't run _all_ of QlikCloud's microservices (yet), but it runs a few, and more are coming
  - Prometheus needs to be able to run with enough CPU/RAM headroom, but not get in the way of other services
- We use service resource _limits_ and _reservations_ - Prometheus gets a 3GB RAM reservation (and limit), and 25% of a CPU reserved, with 150% CPU limit

---

## Where it all started to go wrong...

- We started to notice unusually high `iowait` CPU usage on one of the Swarm Managers in our `stage` (pre-production) Metrics Aggregation Swarm
- Eventually, Docker for AWS's ELB health-check failed (probably due to I/O contention), terminating it and bringing up a new host in its place
  - Swarm considered there to be 3/4 healthy managers (i.e. we still have quorum, but barely)
- I decided _not_ to manually remove the dead manager from the Swarm node list, so I could see what happened _(this was a mistake)_
- Not long after, the high-`iowait` condition started occurring on another manager, and the health-check failure (followed by termination) re-occurred
  - Swarm now considered there to be 2/4 healthy managers. The new manager never managed to join in time.
- _Quorum loss!_ The Swarm was now down.
- I restored the swarm with the `docker swarm init --force-new-cluster ...` command, which will use existing state if it's available (a potential life-saver!)

---

## It gets worse!

- While restoring the swarm, I noticed that Prometheus was scheduled on a Manager node
  - running potentially heavy services is really not a great idea on Swarm Manager nodes
  - we'd forced the target Prometheus instances to not be scheduled on managers, and I totally missed the federator instances _(d'oh!)_
  - simple to fix with Placement Constraints (`--constraint 'node.role != manager'`)
- After applying the constraint, the manager nodes were now totally stable
- _but_ the worker where Prometheus got rescheduled was now experiencing the same `iowait` conditions
- Because it stores samples in-memory, we were still able to query Prometheus about its condition
  - the "urgency score" was hovering around 100% (putting Prometheus into "rushed mode")
  - over the 24hrs prior to the outage, the "chunks to persist" metric had climbed from ~48k to ~480k
- So, Prometheus itself seems to have caused the high `iowait`!

---

## A poor decision comes back to haunt me

- Running an `ls` in Prometheus' data volume took an extremely long time
  - `docker run --rm -v prometheusVol:/data:ro alpine:3.6 ls`
  - this immediately made me suspect EFS...
- If you've been using AWS for a while, you might know that EFS volumes really only _start_ to perform well on _large_ amounts of data (1TB+)
  - but based on performance before this incident, our ~50GB EFS volume seemed OK
  - EFS uses "Burst Credits" to determine throughput
    - burst credits accumulate when you're idle
    - while you have burst credits, you deplete them by reading/writing data at higher speeds (up to 100MB/s)
    - when burst credits hit 0, you're throttled back to the baseline speed (50KB/s per GB of data stored in EFS)
  - 50GB * 50KB/s = 2.5MB/s baseline speed (_this is a really slow filesystem!_)

---

- The last piece of the puzzle came into place when I read this:
  > _Newly created file systems begin with an initial credit balance of 2.1 TiB_
  >
  > _(http://docs.aws.amazon.com/efs/latest/ug/performance.html#bursting)_
- AWS CloudWatch's metrics proved this theory:
.center[<img src="./images/burstcredits.png" style="height: 360px"/>]

---

## 🎉 It's still broken, though

- So as not to lose all our data _just yet_, we "fixed" the problem by throwing more in the EFS volume:
  - `docker run --rm -v make_efs_fast_again:/data alpine:3.6 dd if=/dev/zero of=/data/garbage_data bs=1M count=204800`
  - ran overnight and by morning Prometheus was performing well again

---

## Lessons learned!

1. Don't run your workload on manager nodes.
2. Only use EFS if you have huge amounts of data!
3. Otherwise, use EBS (Cloudstor "relocatable") - it's cheaper too
4. _Definitely_ monitor your services and infrastructure!

---

## References/Links

- [Prometheus](https://prometheus.io)
- [Scaling and Federating Prometheus (blog post)](https://www.robustperception.io/scaling-and-federating-prometheus/)
- [Grafana](https://grafana.com)
- [Docker for AWS](https://docs.docker.com/docker-for-aws/)
- [Docker for AWS - Cloudstor docs](https://docs.docker.com/docker-for-aws/persistent-data-volumes/)

---

## Thank you!

### Any Questions?
