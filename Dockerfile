FROM golang
MAINTAINER Igor Shishkin <me@teran.ru>

ADD . /go
RUN make predependecies dependencies build-linux-amd64

FROM alpine
MAINTAINER Igor Shishkin <me@teran.ru>

RUN apk add --update --no-cache \
  ca-certificates && \
  rm -vf /var/cache/apk/*
COPY --from=0 /go/bin/svcproxy-linux-amd64 /svcproxy

ENTRYPOINT ["/svcproxy"]
