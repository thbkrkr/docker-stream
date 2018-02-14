FROM golang:1.9.4-alpine as builder
ENV GOBIN=/go/bin/ GOPATH=/go
WORKDIR /go/src/github.com/thbkrkr/docker-stream
COPY . /go/src/github.com/thbkrkr/docker-stream
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY entrypoint.sh /entrypoint.sh
COPY --from=builder \
  /go/src/github.com/thbkrkr/docker-stream/docker-stream \
  /usr/local/bin/docker-stream
ENTRYPOINT ["/entrypoint.sh"]