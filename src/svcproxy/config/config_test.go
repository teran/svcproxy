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
        url: http://google.com`)

	cfg, err := parse(configSample)
	s.Require().NoError(err)

}

func (s *ServiceTestSuite) SetupTest() {

}

func TestMyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
