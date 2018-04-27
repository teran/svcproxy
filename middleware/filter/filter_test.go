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

type testCase struct {
	options        map[string]interface{}
	caseIPAddr     string
	caseUAString   string
	expectedStatus int
}

func (s FilterTestSuite) TestAll() {
	// Define test cases
	tcs := []testCase{
		{
			options: map[string]interface{}{
				"name": "filter",
				"rules": []interface{}{
					map[interface{}]interface{}{
						"userAgents": []interface{}{"blah ([0-9]+.[0-9]+)"},
					},
				},
			},
			caseIPAddr:     "127.0.0.1:49000",
			caseUAString:   "SomeBrowser",
			expectedStatus: http.StatusNoContent,
		},
		{
			options: map[string]interface{}{
				"name": "filter",
				"rules": []interface{}{
					map[interface{}]interface{}{
						"userAgents": []interface{}{"blah ([0-9]+.[0-9]+)"},
					},
				},
			},
			caseIPAddr:     "127.0.0.1:49000",
			caseUAString:   "blah 1.0",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			options: map[string]interface{}{
				"name": "filter",
				"rules": []interface{}{
					map[interface{}]interface{}{
						"ips": []interface{}{"127.0.0.1"},
					},
				},
			},
			caseIPAddr:     "127.0.0.1:49000",
			caseUAString:   "blah 1.0",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			options: map[string]interface{}{
				"name": "filter",
				"rules": []interface{}{
					map[interface{}]interface{}{
						"ips": []interface{}{"127.0.0.1"},
					},
				},
			},
			caseIPAddr:     "127.0.0.2:49000",
			caseUAString:   "blah 1.0",
			expectedStatus: http.StatusNoContent,
		},
		{
			options: map[string]interface{}{
				"name": "filter",
				"rules": []interface{}{
					map[interface{}]interface{}{
						"ips":        []interface{}{"127.0.0.1"},
						"userAgents": []interface{}{"blah ([0-9]+.[0-9]+)"},
					},
				},
			},
			caseIPAddr:     "127.0.0.1:49000",
			caseUAString:   "SomeBrowser",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			options: map[string]interface{}{
				"name": "filter",
				"rules": []interface{}{
					map[interface{}]interface{}{
						"ips":        []interface{}{"127.0.0.1"},
						"userAgents": []interface{}{"blah ([0-9]+.[0-9]+)"},
					},
				},
			},
			caseIPAddr:     "127.0.0.2:49000",
			caseUAString:   "blah 1.0",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			options: map[string]interface{}{
				"name": "filter",
				"rules": []interface{}{
					map[interface{}]interface{}{
						"ips":        []interface{}{"127.0.0.1"},
						"userAgents": []interface{}{"blah ([0-9]+.[0-9]+)"},
					},
				},
			},
			caseIPAddr:     "127.0.0.1:49000",
			caseUAString:   "blah 1.0",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			options: map[string]interface{}{
				"name": "filter",
				"rules": []interface{}{
					map[interface{}]interface{}{
						"ips":        []interface{}{"127.0.0.1"},
						"userAgents": []interface{}{"blah ([0-9]+.[0-9]+)"},
					},
				},
			},
			caseIPAddr:     "127.0.0.2:49000",
			caseUAString:   "SomeBrowser",
			expectedStatus: http.StatusNoContent,
		},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	for _, c := range tcs {
		f := NewMiddleware()
		f.SetOptions(c.options)

		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/", nil)
		s.NoError(err)
		r.RemoteAddr = c.caseIPAddr
		r.Header.Set("User-Agent", c.caseUAString)

		f.Middleware(next).ServeHTTP(w, r)

		resp := w.Result()

		s.Equal(c.expectedStatus, resp.StatusCode)
	}
}

func TestFilterTestSuite(t *testing.T) {
	suite.Run(t, new(FilterTestSuite))
}
