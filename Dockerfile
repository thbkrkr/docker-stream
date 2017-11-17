FROM golang:1.9.1 as builder

ENV GOBIN=/go/bin/ GOPATH=/go
WORKDIR /go/src/github.com/thbkrkr/docker-stream
COPY . /go/src/github.com/thbkrkr/docker-stream
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.6
RUN apk --no-cache add ca-certificates
COPY --from=builder \
  /go/src/github.com/thbkrkr/docker-stream/docker-stream \
  /usr/local/bin/docker-stream
ENTRYPOINT ["docker-stream"]