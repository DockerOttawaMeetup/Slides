# DockerCon Recap Demo

_Much of this was inspired by the materials at http://container.training, especially
from the Advanced Orchestration course._

- Bring up 5 nodes with `docker-machine`!
  ```bash
  export AWS_INSTANCE_TYPE=m3.medium
  export AWS_TAGS=Role,docker-demo
  export AWS_SUBNET_ID=subnet-1234abcd
  export AWS_ZONE=d
  export AWS_VPC_ID=vpc-1234abcd
  export AWS_AMI=ami-1234abcd

  export ENGINE_INSTALL_URL=https://experimental.docker.com

  # Set up 5 VMs
  for n in `seq 1 5`; do
    docker-machine create -d amazonec2 --engine-install-url=$ENGINE_INSTALL_URL ec2node${n} &
  done

  # Will need to alter the `docker-machine` security group to allow traffic
  # on port 2377, and optionally port 8000 from the outside world.
  ```
- List their IPs (useful to paste this into `/etc/hosts` for reference later)
  ```bash
  for n in `seq 1 5`; do
    echo -ne "`docker-machine ip ec2node${n}`\tec2node${n}\n"
  done
  ```
- Initialize the swarm with node 1 as manager
  ```bash
  eval $(docker-machine env ec2node1)
  docker swarm init
  ```
- Grab the private IP for node 1
  ```bash
  export NODE1_IP=`docker-machine ssh ec2node1 curl -s http://169.254.169.254/latest/meta-data/local-ipv4`
  ```
- Join the cluster!
  ```bash
  # this one auto-joins
  eval $(docker-machine env ec2node2)
  docker swarm join $NODE1_IP:2377

  # let's set a secret!
  eval $(docker-machine env ec2node1)
  docker swarm update --secret SwarmRulez

  # Now join node 3
  eval $(docker-machine env ec2node3)
  docker swarm join $NODE1_IP:2377
  # uh oh! denied!
  # let's provide the password on join instead
  docker swarm join --secret SwarmRulez $NODE1_IP:2377

  # what about auto-acceptance?
  eval $(docker-machine env ec2node1)
  docker swarm update --auto-accept none

  eval $(docker-machine env ec2node4)
  docker swarm join --secret SwarmRulez $NODE1_IP:2377

  eval $(docker-machine env ec2node1)
  docker node accept XXX

  # Cool, but that can get annoying... Let's turn auto-accept back on for workers
  docker swarm update --auto-accept worker
  docker swarm update --secret ""

  # Finally join node 5
  eval $(docker-machine env ec2node5)
  docker swarm join $NODE1_IP:2377
  ```
- Let's go HA now
  ```bash
  # List the nodes
  eval $(docker-machine env ec2node1)
  docker node ls

  # Now promote 2 and 3
  docker node promote XXX YYY

  # Now if node 1 dies, 2 or 3 will take over as leader!
  ```
- deploy a simple service
  ```bash
  # let's run this from another manager, because!
  eval $(docker-machine env ec2node2)
  docker service create --name pinger alpine ping 8.8.8.8

  # show where it was deployed
  docker service tasks pinger
  ```
- Scaling!
  ```bash
  docker service update pinger --replicas=10
  docker service update pinger --replicas=1
  ```
- Now for something _real!_
  ```bash
  # create an overlay network
  docker network create --driver overlay dockercoins
  docker network ls
  
  docker service create --network dockercoins --name redis redis
  docker service create --network dockercoins --name hasher hairyhenderson/dockercoins_hasher:0.1
  docker service create --network dockercoins --name rng hairyhenderson/dockercoins_rng:0.1
  docker service create --network dockercoins --name worker hairyhenderson/dockercoins_worker:0.1
  # this one needs to publish a port
  docker service create --network dockercoins --name webui \
    --publish 8000:80 hairyhenderson/dockercoins_webui:0.1
  ```
- Check this out at http://ec2node1:8000
  - Now check it out at http://ec2node2:8000 and others... see how the port is
    forwarded, even though it's only running on one node
- Let's scale up!
  ```bash
  docker service update worker --replicas 10
  docker service update rng --replicas 10
  ```
- Time for a rolling update
  ```bash
  docker service update rng --update-parallelism 2 --update-delay 5s
  docker service update rng --image hairyhenderson/dockercoins_rng:0.2
  ```

_Fin._
