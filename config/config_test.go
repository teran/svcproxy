package config

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestConfig() {
	expCfg := &Config{
		Listener: Listener{
			DebugAddr: ":8081",
			HTTPAddr:  ":8080",
			HTTPSAddr: ":8443",
			Frontend: ListenerFrontend{
				IdleTimeout:       5 * time.Second,
				ReadHeaderTimeout: 3 * time.Second,
				ReadTimeout:       5 * time.Second,
				WriteTimeout:      10 * time.Second,
			},
			Backend: ListenerBackend{
				DualStack:             true,
				Timeout:               10 * time.Second,
				KeepAlive:             30 * time.Second,
				ExpectContinueTimeout: 5 * time.Second,
				IdleConnTimeout:       10 * time.Second,
				MaxIdleConns:          10,
				ResponseHeaderTimeout: 10 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
			},
			Middlewares: []map[string]interface{}{
				{
					"name": "filter",
					"rules": []interface{}{
						map[interface{}]interface{}{
							"allowFrom": []interface{}{
								"127.0.0.1/32",
								"::1",
							},
							"denyFrom": []interface{}{
								"127.0.0.2/32",
							},
							"denyUserAgents": []interface{}{
								"blah (Mozilla 5.0)",
							},
						},
					},
				},
				{
					"name": "logging",
				},
				{
					"name": "metrics",
				},
				{
					"name":  "gzip",
					"level": 4,
				},
			},
		},
		Logger: Logger{
			Formatter: "text",
			Level:     "debug",
		},
		Autocert: Autocert{
			Email:        "test@example.com",
			DirectoryURL: "https://acme-v01.api.letsencrypt.org/directory",
			Cache: AutocertCache{
				Backend: "sql",
				BackendOptions: map[string]string{
					"driver":        "mysql",
					"dsn":           "root@tcp(127.0.0.1:3306)/svcproxy?parseTime=true",
					"encryptionKey": "testkey",
					"usePrecaching": "false",
				},
			},
		},
		Services: []Service{
			{
				Frontend: ServiceFrontend{
					FQDN:        []string{"myservice.local", "www.myservice.local"},
					HTTPHandler: "proxy",
					ResponseHTTPHeaders: map[string]string{
						"Strict-Transport-Security": "max-age=31536000",
					},
				},
				Backend: ServiceBackend{
					URL: "http://localhost:8082",
					RequestHTTPHeaders: map[string]string{
						"Host": "example.com",
					},
				},
				Authentication: ServiceAuthentication{
					Method: "BasicAuth",
					Options: map[string]string{
						"backend": "htpasswd",
						"file":    "examples/config/simple/htpasswd",
					},
				},
			},
		},
	}

	configSample, err := ioutil.ReadFile("../examples/config/simple/config.yaml")
	s.Require().NoError(err)

	cfg, err := parse(configSample)
	s.Require().NoError(err)
	s.Require().Equal(expCfg, cfg)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
