package logging

import (
	"log"
	"net"
	"net/http"
	"time"
)

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

// Middleware wraps Handler to log it's request/response metrics
// such as response HTTP status, payload length, time spent.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := ResponseWriterWithStatus{ResponseWriter: w}
		start := time.Now()

		next.ServeHTTP(&rw, r)

		elapsed := time.Since(start)

		remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Error parsing RemoteAddr string: %s", err)
			return
		}

		log.Printf(
			`host="%s" remote_addr="%s" proxy={forwarded_for="%s" forwarded_proto="%s" forwarded_host="%s" real_ip="%s"} method="%s" request_uri="%s" status_code=%d referer="%s" user_agent="%s" duration=%f request_length=%d`,
			r.Host,
			remoteAddr,
			r.Header.Get("X-Forwarded-For"),
			r.Header.Get("X-Forwarded-Proto"),
			r.Header.Get("X-Forwarded-Host"),
			r.Header.Get("X-Real-IP"),
			r.Method,
			r.RequestURI,
			rw.Status,
			r.Referer(),
			r.UserAgent(),
			elapsed.Seconds(),
			r.ContentLength,
		)
	})
}
