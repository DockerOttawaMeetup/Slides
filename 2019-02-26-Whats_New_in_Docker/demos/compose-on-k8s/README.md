# Deploying and using Compose on Kubernetes

https://github.com/docker/compose-on-kubernetes

## Deploying Compose on Kubernetes

Instructions: https://github.com/docker/compose-on-kubernetes/blob/master/docs/install-on-gke.md

### Pre-requisites

1. make sure `helm` is installed
2. Install `compose-on-kubernetes` command:
    - `curl -sSLo ~/bin/compose-on-kubernetes https://github.com/docker/compose-on-kubernetes/releases/download/v0.4.19/installer-darwin`

### Deploy steps

1. `kubectl create namespace compose`
2. Install `etcd`
    - `helm install --name etcd-operator stable/etcd-operator --namespace compose`
    - `kubectl apply -f compose-etcd.yaml`
3. `compose-on-kubernetes -namespace=compose -etcd-servers=http://compose-etcd-client:2379 -tag=v0.4.19`

## Using it!

_(note: for GKE you need a custom docker cli until 19.03 is out)_

```console
$ docker stack deploy --orchestrator=kubernetes -c docker-compose.yml hellokube
```

## Tear-down

```
docker stack rm hellokube
kubectl delete ns compose
```
