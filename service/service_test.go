package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/teran/svcproxy/authentication/factory"
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

	svc, err := NewService()
	s.Require().NoError(err)

	f, err := NewFrontend("test.local", "proxy", map[string]string{"X-Blah": "header-value"})
	s.Require().NoError(err)

	b, err := NewBackend(testsrv.URL)
	s.Require().NoError(err)

	p, err := NewProxy(f, b, nil, nil)
	s.Require().NoError(err)

	svc.AddProxy(p)

	r, err := http.NewRequest("POST", "http://test.local/blah", nil)
	s.Require().NoError(err)

	w := httptest.NewRecorder()

	svc.ServeHTTP(w, r)

	result := w.Result()
	s.Equal(http.StatusNoContent, result.StatusCode)
	s.Equal("header-value", result.Header.Get("X-Blah"))
}

func (s *ServiceTestSuite) TestServiceWithAuthenticator() {
	testsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal("svcproxy", r.Header.Get("X-Proxy-App"))
		s.Equal("/blah", r.URL.Path)
		s.Equal("POST", r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer testsrv.Close()

	svc, err := NewService()
	s.Require().NoError(err)

	f, err := NewFrontend("test.local", "proxy", nil)
	s.Require().NoError(err)

	b, err := NewBackend(testsrv.URL)
	s.Require().NoError(err)

	baOptions := map[string]string{
		"backend": "htpasswd",
		"file":    "../examples/config/simple/htpasswd",
	}

	a, err := factory.NewAuthenticator("BasicAuth", baOptions)
	s.Require().NoError(err)

	p, err := NewProxy(f, b, a, nil)
	s.Require().NoError(err)

	svc.AddProxy(p)

	// Request without authentication
	r, err := http.NewRequest("POST", "http://test.local/blah", nil)
	s.Require().NoError(err)

	w := httptest.NewRecorder()

	svc.ServeHTTP(w, r)

	result := w.Result()
	s.Equal(http.StatusUnauthorized, result.StatusCode)

	// Request with authentication
	r, err = http.NewRequest("POST", "http://test.local/blah", nil)
	s.Require().NoError(err)

	r.SetBasicAuth("testuser", "test")

	w = httptest.NewRecorder()

	svc.ServeHTTP(w, r)

	result = w.Result()

	s.Equal(http.StatusNoContent, result.StatusCode)
}

func (s *ServiceTestSuite) TestRedirect() {
	svc, err := NewService()
	s.Require().NoError(err)

	f, err := NewFrontend("test.local", "redirect", map[string]string{"X-Blah": "header-value"})
	s.Require().NoError(err)

	b, err := NewBackend("http://localhost")
	s.Require().NoError(err)

	p, err := NewProxy(f, b, nil, nil)
	s.Require().NoError(err)

	svc.AddProxy(p)

	r, err := http.NewRequest("POST", "http://test.local/blah?getVar=getValue&getVar2=getValue2", nil)
	s.Require().NoError(err)

	w := httptest.NewRecorder()

	svc.ServeHTTP(w, r)

	result := w.Result()
	s.Equal(http.StatusFound, result.StatusCode)
	s.Equal("", result.Header.Get("X-Blah"))
	s.Equal("https://test.local/blah?getVar=getValue&getVar2=getValue2", result.Header.Get("Location"))
}

func TestMyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
