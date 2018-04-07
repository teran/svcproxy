package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
}

func (s *ServiceTestSuite) TestConfig() {
	configSample := []byte(`---
  listener:
    httpAddr: :80
    httpsAddr: :443
  autocert:
    cache:
      backend: sql
      backendOptions:
        dsn: root:passwd@tcp(127.0.0.1)/svcproxy
  services:
    - frontend:
        fqdn: myservice.local
      backend:
        url: http://example.com`)

	cfg, err := parse(configSample)
	s.Require().NoError(err)
	s.Equal(":80", cfg.Listener.HTTPAddr)
	s.Equal(":443", cfg.Listener.HTTPSAddr)
	s.Equal("sql", cfg.AutoCert.Cache.Backend)
	s.Equal("root:passwd@tcp(127.0.0.1)/svcproxy", cfg.AutoCert.Cache.BackendOptions["dsn"])
	s.Equal("myservice.local", cfg.Services[0].Frontend.FQDN)
	s.Equal("http://example.com", cfg.Services[0].Backend.URL)
}

func (s *ServiceTestSuite) SetupTest() {

}

func TestMyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
