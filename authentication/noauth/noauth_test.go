package noauth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NoAuthTestSuite struct {
	suite.Suite
}

func (s *NoAuthTestSuite) TestNoAuth() {
	r, err := http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)

	na := &NoAuth{}

	isa := na.IsAuthenticated(r)
	s.Require().True(isa)

	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusCreated)

	na.Authenticate(w, r)

	resp := w.Result()

	s.Equal(http.StatusCreated, resp.StatusCode)
}

func TestNoAuthTestSuite(t *testing.T) {
	suite.Run(t, new(NoAuthTestSuite))
}
