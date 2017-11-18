# Docker stream

Push docker container events in Kafka.

```
docker run -ti \
  -e B=${B} \
  -e U=${U} \
  -e P=${P} \
  -e T=${T} \
  -e HOSTNAME=$$(hostname -f)
  -v /var/run/docker.sock:/var/run/docker.sock \
  krkr/docker-stream
```
