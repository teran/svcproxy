package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics type
type Metrics struct {
	httpResponsesTotal    *prometheus.CounterVec
	httpResponseLatencies *prometheus.HistogramVec
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

// NewMetricsMiddleware returns new Middleware
func NewMetricsMiddleware() *Metrics {
	m := Metrics{}

	m.httpResponsesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_responses_total",
			Help: "The count of http responses issued, classified by code, host and method.",
		},
		[]string{"code", "host", "method"},
	)

	m.httpResponseLatencies = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_latencies",
			Help: "The time of http responses issued, classified by code, host and method.",
		},
		[]string{"code", "host", "method"},
	)

	prometheus.MustRegister(m.httpResponsesTotal)
	prometheus.MustRegister(m.httpResponseLatencies)

	return &m
}

// Middleware wraps Handler to obtain metrics
func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := ResponseWriterWithStatus{ResponseWriter: w}
		start := time.Now()

		next.ServeHTTP(&rw, r)

		elapsed := time.Since(start)
		msElapsed := elapsed / time.Millisecond

		m.httpResponsesTotal.WithLabelValues(strconv.Itoa(rw.Status), r.Host, r.Method).Inc()
		m.httpResponseLatencies.WithLabelValues(strconv.Itoa(rw.Status), r.Host, r.Method).Observe(float64(msElapsed))
	})
}
