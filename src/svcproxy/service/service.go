package service

import (
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

var _ Service = &Svc{}
var _ http.Handler = &Svc{}

// Svc implement service
type Svc struct {
	proxies map[string]*httputil.ReverseProxy
}

// NewService returns new service instance
func NewService() (*Svc, error) {
	return &Svc{
		proxies: make(map[string]*httputil.ReverseProxy),
	}, nil
}

// AddProxy adds proxy to the service
func (s *Svc) AddProxy(p *Proxy) error {
	s.proxies[p.Frontend.FQDN] = NewReverseProxy(p)
	return nil
}

func (s *Svc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, ok := s.proxies[r.Host]
	if !ok {
		http.NotFound(w, r)
		return
	}

	p.ServeHTTP(w, r)
}

// NewReverseProxy returns httputil.ReverseProxy object
func NewReverseProxy(p *Proxy) *httputil.ReverseProxy {
	director := func(r *http.Request) {
		r.URL.Scheme = p.Backend.URL.Scheme
		r.URL.Host = p.Backend.URL.Host
		r.URL.Path = singleJoiningSlash(p.Backend.URL.Path, r.URL.Path)

		if p.Backend.RewriteHost {
			r.Host = p.Backend.URL.Host
		}

		if p.Backend.URL.RawQuery == "" || r.URL.RawQuery == "" {
			r.URL.RawQuery = p.Backend.URL.RawQuery + r.URL.RawQuery
		} else {
			r.URL.RawQuery = p.Backend.URL.RawQuery + "&" + r.URL.RawQuery
		}

		remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		if remoteIP == "" {
			remoteIP = "0.0.0.0"
		}

		r.Header.Set("X-Forwarded-For", remoteIP)
		r.Header.Set("X-Real-IP", remoteIP)
		r.Header.Set("X-Forwarded-Proto", "https")
	}

	return &httputil.ReverseProxy{Director: director}
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
