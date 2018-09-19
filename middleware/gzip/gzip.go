package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/teran/svcproxy/middleware/types"
)

var _ types.Middleware = &Gzip{}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Gzip middleware type
type Gzip struct{}

// NewMiddleware returns new Gzip middleware instance
func NewMiddleware() *Gzip {
	return &Gzip{}
}

// SetOptions sets passed options for middleware at startup time(i.e. Chaining procedure)
func (g *Gzip) SetOptions(_ map[string]interface{}) {}

// Middleware wraps handler into GZIP content encoding
func (g *Gzip) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)

		defer gz.Close()

		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		next.ServeHTTP(gzr, r)
	})
}
