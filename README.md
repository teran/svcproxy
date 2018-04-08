# svcproxy

[![Build Status](https://travis-ci.org/teran/svcproxy.svg?branch=master)](https://travis-ci.org/teran/svcproxy)
[![Layers size](https://images.microbadger.com/badges/image/teran/svcproxy.svg)](https://hub.docker.com/r/teran/svcproxy/)
[![Recent build commit](https://images.microbadger.com/badges/commit/teran/svcproxy.svg)](https://hub.docker.com/r/teran/svcproxy/)
[![Docker Automated build](https://img.shields.io/docker/automated/teran/svcproxy.svg)](https://hub.docker.com/r/teran/svcproxy/)
![License](https://img.shields.io/github/license/teran/svcproxy.svg)

HTTP app-agnostic proxy allows to gather metrics and automatically issue certificates using ACME based CA, like Let's Encrypt

# Configuration example

svcproxy uses simple YAML configuration files like this working example:
```
---
listener:
  # Which port to listen for HTTP requests
  httpAddr: :8080
  # Which port to listen for HTTPS requests
  httpsAddr: :8443
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
      encryptionKey: testkey
services:
  - frontend:
      # FQDN service is gonna response by
      fqdn: myservice.local
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
 - [ ] Redirect from HTTP to HTTPS(configurable)
 - [ ] HTTPS-only service
 - [ ] Authentication(?)
 - [X] Fix cache tests
