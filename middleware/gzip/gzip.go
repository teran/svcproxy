package gzip

import (
	"compress/gzip"
	"io"
	"log"
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
type Gzip struct {
	Level int
}

// NewMiddleware returns new Gzip middleware instance
func NewMiddleware() *Gzip {
	return &Gzip{}
}

// SetOptions sets passed options for middleware at startup time(i.e. Chaining procedure)
func (g *Gzip) SetOptions(opts map[string]interface{}) {
	levelField, ok := opts["level"]
	if ok {
		level, ok := levelField.(int)
		if !ok {
			log.Fatal("gzip middleware: error verifying gzip level: probably wrong type, must be integer")
		}

		if level < gzip.HuffmanOnly || level > gzip.BestCompression {
			log.Fatalf("gzip middleware: invalid compression level: %d. Must be between %d and %d", level, gzip.HuffmanOnly, gzip.BestCompression)
		}

		g.Level = level
	}
	g.Level = gzip.DefaultCompression
}

// Middleware wraps handler into GZIP content encoding
func (g *Gzip) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz, err := gzip.NewWriterLevel(w, g.Level)
		if err != nil {
			log.Fatalf("Unexpected error while initializing gzip writer: %s", err)
		}

		defer gz.Close()

		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		next.ServeHTTP(gzr, r)
	})
}
