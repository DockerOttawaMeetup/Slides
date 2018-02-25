#!/bin/sh
export repo=hairyhenderson/helloarch

for arch in arm64 amd64 s390x; do
  docker build --build-arg ARCH=${arch} \
    -t ${repo}:${arch} .
done
for GOARM in 6 7; do
  docker build --build-arg ARCH=arm \
    --build-arg GOARM=${GOARM} \
    -t ${repo}:armv${GOARM} .
done
