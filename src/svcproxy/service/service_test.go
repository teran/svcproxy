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
		s.Equal("svcproxy", r.Header.Get("X-Proxy-App"))
		s.Equal("/blah", r.URL.Path)
		s.Equal("POST", r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer testsrv.Close()

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

	r, err := http.NewRequest("POST", "http://test.local/blah", nil)
	s.Require().NoError(err)

	w := httptest.NewRecorder()

	svc.ServeHTTP(w, r)

	result := w.Result()

	s.Equal(http.StatusNoContent, result.StatusCode)
}

func (s *ServiceTestSuite) SetupTest() {

}

func TestMyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
