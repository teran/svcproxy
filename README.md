# svcproxy

[![Go Report](https://goreportcard.com/badge/github.com/teran/svcproxy)](https://goreportcard.com/report/github.com/teran/svcproxy)
[![Build Status](https://travis-ci.org/teran/svcproxy.svg?branch=master)](https://travis-ci.org/teran/svcproxy)
[![Layers size](https://images.microbadger.com/badges/image/teran/svcproxy.svg)](https://hub.docker.com/r/teran/svcproxy/)
[![Recent build commit](https://images.microbadger.com/badges/commit/teran/svcproxy.svg)](https://hub.docker.com/r/teran/svcproxy/)
[![Docker Automated build](https://img.shields.io/docker/automated/teran/svcproxy.svg)](https://hub.docker.com/r/teran/svcproxy/)
[![License](https://img.shields.io/github/license/teran/svcproxy.svg)](https://github.com/teran/svcproxy/blob/master/LICENSE)

HTTP app-agnostic reverse proxy allows to gather metrics and automatically issue certificates using ACME based CA, like Let's Encrypt

# Configuration example

svcproxy uses simple YAML configuration files like this working example:
```
---
listener:
  # Which address to listen for debug handlers
  # svcproxy will setup handlers for pprof, metrics, tracing
  # on that address.
  # WARNING: this port should never been open to wild Internet!
  debugAddr: :8081
  # Which address to listen for HTTP requests
  httpAddr: :8080
  # Which address to listen for HTTPS requests
  httpsAddr: :8443
  # Middlewares list to apply to each request
  # Available options:
  # - logging
  # - metrics
  # NOTE: amount of middlewares could affect performance and
  #       increase response time.
  middlewares:
    - logging
    - metrics
autocert:
  cache:
    # Cache backend to use
    # Currently available:
    # - sql
    backend: sql
    backendOptions:
      # Driver to use by backend
      # Currently avaialble:
      # - mysql
      # - postgres
      driver: mysql
      # DSN(Data Source Name) to be passed to driver
      dsn: root@tcp(127.0.0.1:3306)/svcproxy
      # PSK(Pre-shared key) to encrypt/decrypt cached data
      # If not set or empty string cache will be used without encryption
      encryptionKey: testkey
services:
  - frontend:
      # FQDN service is gonna response by
      fqdn:
        - myservice.local
        - www.myservice.local
      # What svcproxy should do with requests on HTTP port
      # avaialble options:
      # - "proxy" to work on both of HTTP and HTTPS
      # - "redirect" to redirect requests from HTTP to HTTPS
      # - "reject" to reject any requests to HTTP(except ACME challenges) with 404
      httpHandler: proxy
      # HTTP Headers to send with response
      # Usually usefull for HSTS, CORS, etc.
      responseHTTPHeaders:
        Strict-Transport-Security: "max-age=31536000"
    backend:
      # Service backend to handle requests behind proxy
      url: http://localhost:8082
```

Some options could be passed as Environment variables:
 * `CONFIG_PATH` - path to YAML configuration file in file system

# Builds

Automatic builds are available on DockerHub:
```
docker pull teran/svcproxy
```

# TODO
 - [X] Redirect from HTTP to HTTPS(configurable)
 - [X] HTTPS-only service
 - [X] Fix cache tests
 - [X] Multiple names for proxy(aliases)
 - [ ] Autocert SQL cache to cache certificates in memory(reduce amount of SELECT's)
 - [ ] Autocert cache for Redis or Mongo (?)
 - [ ] Authentication(?)
 - [ ] Tracing(?)
