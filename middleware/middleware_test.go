package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
}

func (s *MiddlewareTestSuite) TestAll() {
	mdlwrs := []map[string]interface{}{
		{
			"name":  "gzip",
			"level": 4,
		},
		{
			"name": "logging",
		},
		{
			"name": "metrics",
		},
		{
			"name": "filter",
			"rules": []map[string][]string{
				map[string][]string{
					"allowFrom": []string{
						"127.0.0.1/32",
					},
					"denyFrom": []string{
						"127.0.0.2/32",
					},
					"denyUserAgents": []string{
						"blah (Mozilla 5.0)",
					},
				},
			},
		},
	}
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusMultiStatus), http.StatusMultiStatus)
	})
	hndlr, err := Chain(f, mdlwrs...)
	s.Require().NoError(err)
	s.Require().NotNil(hndlr)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, &MiddlewareTestSuite{})
}
