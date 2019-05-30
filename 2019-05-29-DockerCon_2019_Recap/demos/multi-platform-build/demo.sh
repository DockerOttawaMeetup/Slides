#!/bin/sh
#
# multi-platform build demo, using gomplate

AMD_HOST=ssh://hairyhenderson@radon.local
ARM_HOST=ssh://hairyhenderson@tin.local
# AMD_HOST=unix:///var/run/docker.sock
ARM64_HOST=ssh://ubuntu@18.212.222.14

# first we need some builders that support our multi-platform build
docker buildx create \
    --name multi-builder \
    --node arm64-node \
    --platform linux/arm64 \
    $ARM64_HOST
docker buildx create \
    --append \
    --name multi-builder \
    --node arm-node \
    --platform linux/arm/v7,linux/arm/v6 \
    $ARM_HOST
docker buildx create \
    --append \
    --name multi-builder \
    --node amd-node \
    --platform linux/amd64 \
    $AMD_HOST
# docker buildx create \
#     --append \
#     --name multi-builder \
#     --node local-node \
#     --platform linux/amd64

# can also just build arm locally too, since Docker Desktop includes QEMU configured for that:
# docker buildx create \
#     --name multi-builder \
#     --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6

# we can see the different builders we have
docker buildx ls

# ...and use the one we want
docker buildx use multi-builder

# then we can build it!
git clone git@github.com:hairyhenderson/dockerfiles.git
cd dockerfiles/figlet
docker buildx build \
    --platform linux/amd64,linux/arm64,linux/arm \
    -t hairyhenderson/figlet:multi \
    --push \
    .


# clean-up
rm -rf dockerfiles
docker buildx rm multi-builder
