package gzip

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GZipMiddlewareTestSuite struct {
	suite.Suite
}

var _ http.Handler = &testHandler{}

var handlerContent = "this is a default handler content"

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, handlerContent)
}

type testCase struct {
	level                    int
	acceptEncoding           string
	expContentEncodingHeader string
	expContentLengthHeader   string
	expErrorOnSetConfig      bool
}

func (s *GZipMiddlewareTestSuite) TestAll() {
	tcs := func() []testCase {
		tt := []testCase{}
		for i := 0; i <= 9; i++ {
			tt = append(tt, testCase{
				level:                    i,
				acceptEncoding:           "gzip",
				expContentEncodingHeader: "gzip",
				expContentLengthHeader:   "",
				expErrorOnSetConfig: func(i int) bool {
					if 0 <= i && i <= 9 {
						return false
					}
					return true
				}(i),
			})
		}

		return tt
	}()

	for _, tc := range tcs {
		g := NewMiddleware()
		err := g.SetConfig(&GzipConfig{
			Level: tc.level,
		})
		if tc.expErrorOnSetConfig {
			s.Require().Error(err)
		} else {
			s.Require().NoError(err)
		}

		s.Require().Equal(tc.level, g.(*Gzip).Level)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		s.Require().NoError(err)

		req.Header.Set("Accept-Encoding", tc.acceptEncoding)

		testsrv := &testHandler{}
		g.Middleware(testsrv).ServeHTTP(w, req)

		result := w.Result()
		s.Require().Equal(tc.expContentEncodingHeader, result.Header.Get("Content-Encoding"))
		s.Require().Equal(tc.expContentLengthHeader, result.Header.Get("Content-Length"))

		gzreader, err := gzip.NewReader(result.Body)
		s.Require().NoError(err)
		defer result.Body.Close()

		uncompressedBody, err := ioutil.ReadAll(gzreader)
		s.Require().NoError(err)
		s.Require().Equal(handlerContent, string(uncompressedBody))
	}
}

func TestGZipMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, &GZipMiddlewareTestSuite{})
}
