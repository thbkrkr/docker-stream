#!/bin/sh -eu

case ${1:-""} in
  up)
    echo '
    docker pull krkr/docker-stream
    docker run --name docker-stream -d \
      -e HOSTNAME=${DOCKER_MACHINE_NAME} \
      -e B=${B} \
      -e U=${U} \
      -e P=${P} \
      -e T=${T} \
      -v /var/run/docker.sock:/var/run/docker.sock \
      krkr/docker-stream'
    ;;
  *)
    exec docker-stream
    ;;
esac