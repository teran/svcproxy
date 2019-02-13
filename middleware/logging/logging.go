package logging

import (
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/teran/svcproxy/middleware/types"
)

var _ types.Middleware = (*Logging)(nil)

// Logging middleware type
type Logging struct{}

// NewMiddleware returns new Middleware instance
func NewMiddleware() *Logging {
	return &Logging{}
}

// SetOptions sets passed options for middleware at startup time(i.e. Chaining procedure)
func (l *Logging) SetOptions(_ map[string]interface{}) {}

// Middleware wraps Handler to log it's request/response metrics
// such as response HTTP status, payload length, time spent.
func (l *Logging) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := ResponseWriterWithStatus{ResponseWriter: w}
		start := time.Now()

		next.ServeHTTP(&rw, r)

		elapsed := time.Since(start)

		remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Warnf("Error parsing RemoteAddr string: %s", err)
			return
		}

		log.WithFields(log.Fields{
			"host":            strings.ToLower(r.Host),
			"remote_addr":     remoteAddr,
			"forwarded_for":   r.Header.Get("X-Forwarded-For"),
			"forwarded_proto": r.Header.Get("X-Forwarded-Proto"),
			"forwarded_host":  r.Header.Get("X-Forwarded-Host"),
			"real_ip":         r.Header.Get("X-Real-IP"),
			"method":          r.Method,
			"request_uri":     r.RequestURI,
			"status_code":     rw.Status,
			"referer":         r.Referer(),
			"user_agent":      r.UserAgent(),
			"duration":        elapsed.Seconds(),
			"request_length":  r.ContentLength,
		}).Info("Request handled")
	})
}

// ResponseWriterWithStatus implements adding status code to ResponseWriter object
type ResponseWriterWithStatus struct {
	http.ResponseWriter
	Status int
}

// WriteHeader reimplements WriteHeader() to fill status automatically
func (rw *ResponseWriterWithStatus) WriteHeader(status int) {
	rw.Status = status
	rw.ResponseWriter.WriteHeader(status)
}
