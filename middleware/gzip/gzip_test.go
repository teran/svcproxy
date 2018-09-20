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

type GZipMiddlewareSuite struct {
	suite.Suite
}

var _ http.Handler = &testHandler{}

var handlerContent = "this is a default handler content"

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, handlerContent)
}

func (s *GZipMiddlewareSuite) TestGZipMiddleware() {
	g := NewMiddleware()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)

	req.Header.Set("Accept-Encoding", "gzip")

	testsrv := &testHandler{}
	g.Middleware(testsrv).ServeHTTP(w, req)

	result := w.Result()
	s.Require().Equal("gzip", result.Header.Get("Content-Encoding"))

	gzreader, err := gzip.NewReader(result.Body)
	s.Require().NoError(err)
	defer result.Body.Close()

	uncompressedBody, err := ioutil.ReadAll(gzreader)
	s.Require().NoError(err)
	s.Require().Equal(handlerContent, string(uncompressedBody))
}

func (s *GZipMiddlewareSuite) TestGZipMiddlewareNoEncoding() {
	g := NewMiddleware()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	s.Require().NoError(err)

	testsrv := &testHandler{}
	g.Middleware(testsrv).ServeHTTP(w, req)

	result := w.Result()
	s.Require().Equal("", result.Header.Get("Content-Encoding"))

	uncompressedBody, err := ioutil.ReadAll(result.Body)
	defer result.Body.Close()

	s.Require().NoError(err)
	s.Require().Equal(handlerContent, string(uncompressedBody))
}

func TestGZipMiddlewareSuite(t *testing.T) {
	suite.Run(t, &GZipMiddlewareSuite{})
}
