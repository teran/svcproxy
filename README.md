# svcproxy
HTTP app-agnostic proxy allows to gather metrics and automatically issue certificates using ACME based CA, like Let's Encrypt

# Configuration example

svcproxy uses simple YAML configuration files like this working example:
```
---
listener:
  httpAddr: :8080
  httpsAddr: :8443
autocert:
  cache:
    backend: sql
    backendOptions:
      driver: mysql
      dsn: root@tcp(127.0.0.1:3306)/svcproxy
      encryptionKey: testkey
services:
  - frontend:
      fqdn: myservice.local
    backend:
      url: http://example.com
```

Some options could be passed as Environment variables:
 * `CONFIG_PATH` - path to YAML configuration file in file system
