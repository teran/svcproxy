package config

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestConfig() {
	configSample, err := ioutil.ReadFile("../examples/config/simple/config.yaml")
	s.Require().NoError(err)

	cfg, err := parse(configSample)
	s.Require().NoError(err)
	s.Equal(":8080", cfg.Listener.HTTPAddr)
	s.Equal(":8443", cfg.Listener.HTTPSAddr)
	s.Equal("sql", cfg.Autocert.Cache.Backend)
	s.Equal("mysql", cfg.Autocert.Cache.BackendOptions["driver"])
	s.Equal("root@tcp(127.0.0.1:3306)/svcproxy?parseTime=true", cfg.Autocert.Cache.BackendOptions["dsn"])
	s.Equal([]string{"myservice.local", "www.myservice.local"}, cfg.Services[0].Frontend.FQDN)
	s.Equal("http://localhost:8082", cfg.Services[0].Backend.URL)
}

func (s *ConfigTestSuite) SetupTest() {

}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
