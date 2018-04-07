FROM golang
MAINTAINER Igor Shishkin <me@teran.ru>

ADD . /go
RUN make predependecies dependencies
RUN make build-linux-amd64

FROM scratch
MAINTAINER Igor Shishkin <me@teran.ru>

COPY --from=0 /go/bin/svcproxy-linux-amd64 /svcproxy

ENTRYPOINT ["/svcproxy"]
