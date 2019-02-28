# Demo 1

A simple performance/optimization demo.

We're going to avoid the cache and pull everything, to show the pathological case. Later, we can re-build to show caching at work...

1. Build with the legacy builder
    ```bash
    export DOCKER_BUILDKIT=0
    docker rmi alpine:3.8
    docker image prune
    time docker build -t demo1 .
    ```
    sample result:
    ```
    ...0.08s user 0.06s system 0% cpu 17.307 total
    ```
2. now build with BuildKit and watch it fly!
    ```bash
    export DOCKER_BUILDKIT=1
    docker rmi alpine:3.8
    docker image prune
    docker builder prune
    time docker build -t demo1 .
    ```
    sample result:
    ```
    ...0.09s user 0.07s system 3% cpu 4.399 total
    ```
