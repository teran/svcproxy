package service

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
}

func (s *ServiceTestSuite) TestService() {
	u, err := url.Parse("http://example.com")
	s.Require().NoError(err)

	svc, err := NewService()
	s.Require().NoError(err)

	svc.AddProxy(&Proxy{
		Frontend: &Frontend{
			FQDN: "test.local",
		},
		Backend: &Backend{
			URL: u,
		},
	})
}

func (s *ServiceTestSuite) SetupTest() {

}

func TestMyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
