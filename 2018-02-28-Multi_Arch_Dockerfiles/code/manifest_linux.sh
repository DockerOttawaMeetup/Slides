#!/bin/sh

repo=hairyhenderson/helloarch
docker manifest create --amend ${repo}:latest \
    ${repo}:armv6 \
    ${repo}:armv7 \
    ${repo}:arm64 \
    ${repo}:amd64 \
    ${repo}:s390x

# see https://github.com/opencontainers/image-spec/blob/v1.0.1/image-index.md
# for valid variants
docker manifest annotate ${repo}:latest ${repo}:armv6 \
    --arch arm --variant v6
docker manifest annotate ${repo}:latest ${repo}:armv7 \
    --arch arm --variant v7
docker manifest annotate ${repo}:latest ${repo}:arm64 \
    --arch arm64 --variant v8
docker manifest annotate ${repo}:latest ${repo}:s390x \
    --arch s390x
