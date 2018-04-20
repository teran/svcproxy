package basicauth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"svcproxy/authentication"
)

const testUsername = "gotest"
const testPassword = "testpassword"

type BasicAuthTestSuite struct {
	suite.Suite
	basicAuth authentication.Authenticator
}

func (s *BasicAuthTestSuite) TestAuthenticationNoCredentials() {
	r, err := http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)
	w := httptest.NewRecorder()

	isa := s.basicAuth.IsAuthenticated(r)
	s.False(isa)

	s.basicAuth.Authenticate(w, r)

	resp := w.Result()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	s.Equal(`Basic realm="Restricted area"`, resp.Header.Get("WWW-Authenticate"))
}

func (s *BasicAuthTestSuite) TestAuthenticationWithValidCredentials() {
	r, err := http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)

	r.SetBasicAuth(testUsername, testPassword)
	w := httptest.NewRecorder()

	isa := s.basicAuth.IsAuthenticated(r)
	s.True(isa)

	s.basicAuth.Authenticate(w, r)

	resp := w.Result()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	s.Equal(`Basic realm="Restricted area"`, resp.Header.Get("WWW-Authenticate"))
}

func (s *BasicAuthTestSuite) TestAuthenticationWithInvalidPassword() {
	r, err := http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)

	r.SetBasicAuth(testUsername, "wrongPassword")
	w := httptest.NewRecorder()

	isa := s.basicAuth.IsAuthenticated(r)
	s.False(isa)

	s.basicAuth.Authenticate(w, r)

	resp := w.Result()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	s.Equal(`Basic realm="Restricted area"`, resp.Header.Get("WWW-Authenticate"))
}

func (s *BasicAuthTestSuite) TestAuthenticationWithInvalidUsername() {
	r, err := http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)

	r.SetBasicAuth("wrongUsername", "wrongPassword")
	w := httptest.NewRecorder()

	isa := s.basicAuth.IsAuthenticated(r)
	s.False(isa)

	s.basicAuth.Authenticate(w, r)

	resp := w.Result()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	s.Equal(`Basic realm="Restricted area"`, resp.Header.Get("WWW-Authenticate"))
}

func (s *BasicAuthTestSuite) SetupTest() {
	hash, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	s.Require().NoError(err)

	dir, err := ioutil.TempDir("", "gotestPasswd")
	s.Require().NoError(err)

	passwdFile := path.Join(dir, "passwd")

	err = ioutil.WriteFile(passwdFile, []byte(fmt.Sprintf("%s:%s", testUsername, hash)), 0644)
	s.Require().NoError(err)

	s.basicAuth, err = NewBasicAuth("htpasswd", map[string]string{"file": passwdFile})
	s.Require().NoError(err)
}

func TestMyBasicAuthTestSuite(t *testing.T) {
	suite.Run(t, new(BasicAuthTestSuite))
}
