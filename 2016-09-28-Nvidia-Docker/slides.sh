#!/bin/sh

slide="slide-nvidia-docker"
port=8800

echo "usage: cmd [ port=${port} ]"
echo

slide_path=${USER_WORKSPACE_HOME}/github/nvidia-docker/slides

if [ $# -gt 0 ]; then
    port=${1}
fi

set -x

docker stop ${slide}
docker rm ${slide}
docker run --name ${slide} -d -p ${port}:1948 -v ${slide_path}/images:/usr/src/app/images -v ${slide_path}/slide.md:/usr/src/app/slide.md  muellermich/reveal-md:latest
