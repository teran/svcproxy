# svcproxy

[![Go Report](https://goreportcard.com/badge/github.com/teran/svcproxy)](https://goreportcard.com/report/github.com/teran/svcproxy)
[![Build Status](https://travis-ci.org/teran/svcproxy.svg?branch=master)](https://travis-ci.org/teran/svcproxy)
[![License](https://img.shields.io/github/license/teran/svcproxy.svg)](https://github.com/teran/svcproxy/blob/master/LICENSE)

[![Layers size](https://images.microbadger.com/badges/image/teran/svcproxy.svg)](https://hub.docker.com/r/teran/svcproxy/)
[![Recent build commit](https://images.microbadger.com/badges/commit/teran/svcproxy.svg)](https://hub.docker.com/r/teran/svcproxy/)
[![Docker Automated build](https://img.shields.io/docker/automated/teran/svcproxy.svg)](https://hub.docker.com/r/teran/svcproxy/)
[![GoDoc](https://godoc.org/github.com/teran/svcproxy?status.svg)](https://godoc.org/github.com/teran/svcproxy)

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
  # Frontend global settings
  frontend:
    # idleTimeout is passed as IdleTimeout to http.Server which is described as:
    #   IdleTimeout is the maximum amount of time to wait for the
    #   next request when keep-alives are enabled. If IdleTimeout
    #   is zero, the value of ReadTimeout is used. If both are
    #   zero, ReadHeaderTimeout is used.
    idleTimeout: 5s
    # readHeaderTimeout is passed as ReadHeaderTimeout to http.Server which is
    # described as:
    #   ReadHeaderTimeout is the amount of time allowed to read
    #   request headers. The connection's read deadline is reset
    #   after reading the headers and the Handler can decide what
    #   is considered too slow for the body.
    readHeaderTimeout: 3s
    # readTimeout is passed as ReadTimeout to http.Server which is described as:
    #   ReadTimeout is the maximum duration for reading the entire
    #   request, including the body.
    #
    #   Because ReadTimeout does not let Handlers make per-request
    #   decisions on each request body's acceptable deadline or
    #   upload rate, most users will prefer to use
    #   ReadHeaderTimeout. It is valid to use them both.
    readTimeout: 5s
    # writeTimeout is passed as WriteTimeout to http.Server which is described as:
    #   WriteTimeout is the maximum duration before timing out
    #   writes of the response. It is reset whenever a new
    #   request's header is read. Like ReadTimeout, it does not
    #   let Handlers make decisions on a per-request basis.
    writeTimeout: 10s
  # Backend global settings
  backend:
    # More details about the following options could be found at:
    #   https://golang.org/pkg/net/#Dialer
    dualStack: true
    timeout: 10s
    keepAlive: 30s
    # More details about the following options could be found at:
    #   https://golang.org/pkg/net/http/#Transport
    expectContinueTimeout: 5s
    idleConnTimeout: 10s
    maxIdleConns: 10
    responseHeaderTimeout: 10s
    tlsHandshakeTimeout: 10s
  # Middlewares list to apply to each request passing through HTTPS socket
  # Available options:
  # - filter
  # - logging
  # - metrics
  # NOTE: amount of middlewares could affect performance and
  #       increase response time.
  middlewares:
    - name: filter
      rules:
        - allowFrom:
          - "127.0.0.1/32"
          - "::1"
          denyFrom:
          - "127.0.0.2/32"
          denyUserAgents:
          - "blah (Mozilla 5.0)"
    - name: logging
    - name: metrics
    - name: gzip
logger:
  # Log formatter to use. Available options are: text, json
  formatter: text
  # Log verbosity. Available options are: debug, info, warning, error, fatal, panic
  level: debug
autocert:
  # Email optionally specifies a contact email address.
  # This is used by CAs, such as Let's Encrypt, to notify about problems
  # with issued certificates.
  email: test@example.com
  # CA Directory endpoint URL
  # Could be left empty or not specified to use Let's Encrypt
  # Default: https://acme-v01.api.letsencrypt.org/directory
  directoryURL: "https://acme-v01.api.letsencrypt.org/directory"
  # Local cache settings
  cache:
    # Cache backend to use
    # Currently available:
    # - dir
    # - redis
    # - sql
    # More details about configuration at:
    #   https://github.com/teran/svcproxy/blob/master/autocert/cache/README.md
    backend: sql
    backendOptions:
      # Driver to use by backend
      # Currently avaialble:
      # - mysql
      # - postgres
      driver: mysql
      # DSN(Data Source Name) to be passed to driver
      # NOTE: parseTime option is required for MySQL driver to be true for
      #       migrations engine
      dsn: root@tcp(127.0.0.1:3306)/svcproxy?parseTime=true
      # PSK(Pre-shared key) to encrypt/decrypt cached data
      # If not set or empty string cache will be used without encryption
      encryptionKey: testkey
      # Precache certificates in memory in unencrypted form to make it much-much
      # faster, faster as serve from memory. default = false.
      # Supported in all of the available cache backends.
      # WARNING: this could decrease security of the certificates
      # WARNING: this will decrease security and could cause certificates leaks
      #          in case of core dumps turned on
      usePrecaching: false
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
      # Request headers passed to backend
      requestHTTPHeaders:
        Host: example.com
    # Authnticator to use for current proxy
    # Currently available:
    # - BasicAuth
    # - NoAuth (default)
    authentication:
      method: BasicAuth
      # Options to pass to authenticator, normally depends on what is supported
      # by particular authenticator
      # For BasicAuth supported options:
      # - backend (backend to use by BasicAuth authenticator)
      # - file(used by htpasswd backend), path to htpasswd file
      options:
        backend: htpasswd
        file: examples/config/simple/htpasswd
```


Some options could be passed as Environment variables:
 * `CONFIG_PATH` - path to YAML configuration file in file system

# Builds

Automatic builds are available on DockerHub:
```
docker pull teran/svcproxy
```

# Authentication
## BasicAuth
### htpasswd backend

htpasswd backend implements simple Basic Auth mechanism via HTTP headers(rfc2617),
using `htpasswd` file as a user database(Bcrypt only is supported).

To generate `htpasswd` file for svcproxy please use the following command:

```
htpasswd -Bc <filename> <username>
```

Please note, `htpasswd` CLI is not vendored with Docker image or in any other way
with svcproxy, but could be easily obtained from packge repositories like `Homebrew`, `ubuntu.archive.com`, etc.

# Using as library

Some parts of svcproxy like [autocert.Cache](https://godoc.org/golang.org/x/crypto/acme/autocert#Cache) implementations are awesome to be used as
libraries in some third-party software. Since all the cache subsystem is layered and implemented as
[autocert.Cache](https://godoc.org/golang.org/x/crypto/acme/autocert#Cache) on all of the layers it could be easily used in any of the following ways:

```
import "github.com/teran/svcproxy/autocert/cache"

....
c, err := cache.NewCacheFactory(.......)
....
```

```
import "github.com/teran/svcproxy/autocert/cache/redis"

....
c, err := redis.NewCache(.......)
....
```

Please note, using higher level `github.com/teran/svcproxy/autocert/cache.Cache` instance
will enable encryption and precaching for your app.

Feel free to take a look at [GoDoc](https://godoc.org/github.com/teran/svcproxy/autocert/cache) for
further details.

# TODO
 - [X] Redirect from HTTP to HTTPS(configurable)
 - [X] HTTPS-only service
 - [X] Fix cache tests
 - [X] Multiple names for proxy(aliases)
 - [X] Autocert SQL cache to cache certificates in memory(reduce amount of SELECT's)
 - [X] Authentication(?)
 - [X] Autocert cache for Redis or Mongo (?)
 - [ ] Tracing(?)
