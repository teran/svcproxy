FROM golang
MAINTAINER Igor Shishkin <me@teran.ru>

ADD . /go
RUN make build-linux-amd64

FROM alpine
MAINTAINER Igor Shishkin <me@teran.ru>

ARG VCS_REF

LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="e.g. https://github.com/teran/svcproxy"

RUN apk add --update --no-cache \
  ca-certificates && \
  rm -vf /var/cache/apk/*
COPY --from=0 /go/bin/svcproxy-linux-amd64 /svcproxy

ENTRYPOINT ["/svcproxy"]
