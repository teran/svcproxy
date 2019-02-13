package gzip

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/teran/svcproxy/middleware/types"
)

var _ types.Middleware = (*Gzip)(nil)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipConfig type
type GzipConfig struct {
	Level int
}

// Unpack implemnts types.MiddlewareConfig interface
func (mc *GzipConfig) Unpack(options map[string]interface{}) error {
	lvl, ok := options["level"].(int)
	if ok {
		mc.Level = lvl
		return nil
	}
	return errors.New("gzip compression level must be defined")
}

// Gzip middleware type
type Gzip struct {
	Level int
}

// NewMiddleware returns new Gzip middleware instance
func NewMiddleware() types.Middleware {
	return &Gzip{}
}

// SetConfig applies config to the middleware
func (g *Gzip) SetConfig(opts types.MiddlewareConfig) error {
	o := opts.(*GzipConfig)
	if o.Level < gzip.HuffmanOnly || o.Level > gzip.BestCompression {
		return fmt.Errorf("gzip middleware: invalid compression level: %d. Must be between %d and %d", o.Level, gzip.HuffmanOnly, gzip.BestCompression)
	}

	g.Level = o.Level

	return nil
}

// Middleware wraps handler into GZIP content encoding
func (g *Gzip) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")

		gz, err := gzip.NewWriterLevel(w, g.Level)
		if err != nil {
			log.Fatalf("Unexpected error while initializing gzip writer: %s", err)
		}

		defer gz.Close()

		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		next.ServeHTTP(gzr, r)
	})
}
