package service

import (
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _ Service = &Svc{}
var _ http.Handler = &Svc{}

// Svc implement service
type Svc struct {
	proxies map[string]*Proxy
}

// NewService returns new service instance
func NewService() (*Svc, error) {
	return &Svc{
		proxies: make(map[string]*Proxy),
	}, nil
}

// AddProxy adds proxy to the service
func (s *Svc) AddProxy(p *Proxy) error {
	s.proxies[p.Frontend.FQDN] = p
	return nil
}

func (s *Svc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, ok := s.proxies[r.Host]
	if !ok {
		http.NotFound(w, r)
		return
	}

	for k, v := range p.Frontend.ResponseHTTPHeaders {
		w.Header().Set(k, v)
	}

	// Handle plain HTTP requests
	if r.TLS == nil {
		switch p.Frontend.HTTPHandler {
		case "reject":
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		case "redirect":
			host, _, err := net.SplitHostPort(r.Host)
			if err != nil {
				// absence of port in address causes error in SplitHostPort()
				// so hope, it's our case :)
				host = r.Host
			}

			redirURL := &url.URL{
				Scheme:   "https",
				Host:     host,
				Path:     r.URL.Path,
				RawQuery: r.URL.RawQuery,
			}

			http.Redirect(w, r, redirURL.String(), http.StatusFound)
			return
		}
	}

	p.proxy.ServeHTTP(w, r)
}

// DebugHandlerFunc implements handlers for debug listener
func (s *Svc) DebugHandlerFunc(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	mux.Handle("/health/metrics", promhttp.Handler())
	mux.Handle("/health/ping", http.HandlerFunc(s.debugPing))

	mux.ServeHTTP(w, r)
}

func (s *Svc) debugPing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
