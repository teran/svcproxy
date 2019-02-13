package metrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/teran/svcproxy/middleware/types"
)

var _ types.Middleware = (*Metrics)(nil)

// Metrics type
type Metrics struct {
	inFlightRequests           prometheus.Gauge
	httpRequestsTotal          *prometheus.CounterVec
	responseDurationSeconds    *prometheus.HistogramVec
	writeHeaderDurationSeconds *prometheus.HistogramVec
	responseSizeBytes          *prometheus.HistogramVec
	requestSizeBytes           *prometheus.HistogramVec
}

// ResponseWriterWithStatus implements adding status code to ResponseWriter object
type ResponseWriterWithStatus struct {
	http.ResponseWriter
	Status                 int
	Written                int64
	ObserveWriteHeaderFunc func(int)
}

// WriteHeader reimplements WriteHeader() to fill status automatically
func (rw *ResponseWriterWithStatus) WriteHeader(status int) {
	rw.Status = status
	rw.ResponseWriter.WriteHeader(status)
	if rw.ObserveWriteHeaderFunc != nil {
		rw.ObserveWriteHeaderFunc(status)
	}
}

func (rw *ResponseWriterWithStatus) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.Written += int64(n)
	return n, err
}

// NewMiddleware returns new Middleware instance
func NewMiddleware() *Metrics {
	m := Metrics{}

	m.inFlightRequests = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	m.httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"host", "code", "method"},
	)

	m.responseDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "A histogram of request latencies.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"host", "code", "method"},
	)

	m.writeHeaderDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_write_header_duration_seconds",
			Help:    "A histogram of time to first write latencies.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"host", "code", "method"},
	)

	m.requestSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "A histogram of request sizes.",
			Buckets: []float64{50, 200, 500, 900, 1500},
		},
		[]string{"host", "code", "method"},
	)

	m.responseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "A histogram of response sizes.",
			Buckets: []float64{50, 200, 500, 900, 1500},
		},
		[]string{"host", "code", "method"},
	)

	prometheus.MustRegister(m.inFlightRequests)
	prometheus.MustRegister(m.httpRequestsTotal)
	prometheus.MustRegister(m.responseDurationSeconds)
	prometheus.MustRegister(m.writeHeaderDurationSeconds)
	prometheus.MustRegister(m.requestSizeBytes)
	prometheus.MustRegister(m.responseSizeBytes)

	return &m
}

// SetOptions sets passed options for middleware at startup time(i.e. Chaining procedure)
func (m *Metrics) SetOptions(_ map[string]interface{}) {}

// Middleware wraps Handler to obtain metrics
func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostName := strings.ToLower(r.Host)
		now := time.Now()
		rw := ResponseWriterWithStatus{
			ResponseWriter: w,
			ObserveWriteHeaderFunc: func(status int) {
				m.writeHeaderDurationSeconds.WithLabelValues(hostName, strconv.Itoa(status), r.Method).Observe(time.Since(now).Seconds())
			},
		}

		m.inFlightRequests.Inc()
		defer m.inFlightRequests.Dec()

		next.ServeHTTP(&rw, r)

		statusCode := strconv.Itoa(rw.Status)

		m.responseDurationSeconds.WithLabelValues(hostName, statusCode, r.Method).Observe(time.Since(now).Seconds())
		m.httpRequestsTotal.WithLabelValues(hostName, statusCode, r.Method).Inc()
		m.requestSizeBytes.WithLabelValues(hostName, statusCode, r.Method).Observe(float64(calculateRequestSize(r)))
		m.responseSizeBytes.WithLabelValues(hostName, statusCode, r.Method).Observe(float64(rw.Written))
	})
}

// Calculate (approximately) request size
func calculateRequestSize(r *http.Request) int {
	size := 0

	size += len(r.Method)
	size += len(r.URL.Path)
	size += len(r.Proto)

	// Add 6 bytes for "Host: "
	size += 6
	size += len(r.Host)

	for header, values := range r.Header {
		size += len(header)
		for _, value := range values {
			// Add 2 bytes for semicolon and space usually present between
			// header name and it's value
			size += 2
			size += len(value)
		}
	}

	if r.ContentLength != -1 {
		size += int(r.ContentLength)
	}

	return size
}
