package filter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FilterTestSuite struct {
	suite.Suite
}

func (s *FilterTestSuite) TestFilter() {
	f := NewMiddleware()
	f.SetOptions(map[string]interface{}{
		"name": "filter",
		"rules": []interface{}{
			map[interface{}]interface{}{
				"userAgents": []interface{}{"blah ([0-9]+.[0-9]+)"},
			},
		},
	})

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Test successful pass
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)
	r.RemoteAddr = "127.0.0.1:4443"

	f.Middleware(next).ServeHTTP(w, r)

	resp := w.Result()

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Test filtered request
	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)

	r.RemoteAddr = "127.0.0.1:4443"
	r.Header.Set("User-Agent", "blah 1.0")

	f.Middleware(next).ServeHTTP(w, r)

	resp = w.Result()

	s.Equal(http.StatusServiceUnavailable, resp.StatusCode)
}

func TestFilterTestSuite(t *testing.T) {
	suite.Run(t, new(FilterTestSuite))
}
