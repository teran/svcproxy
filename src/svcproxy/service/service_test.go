package service

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
}

func (s *ServiceTestSuite) TestService() {
	testsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	}))

	u, err := url.Parse(testsrv.URL)
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
