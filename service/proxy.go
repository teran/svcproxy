package service

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/teran/svcproxy/authentication"
)

// NewProxy creates new Proxy instance
func NewProxy(frontend *Frontend, backend *Backend, authenticator authentication.Authenticator, transport http.RoundTripper, logger *log.Logger) (*Proxy, error) {
	rp := NewReverseProxy(backend, transport)

	if logger != nil {
		rp.ErrorLog = logger
	}

	p := &Proxy{
		Frontend:      frontend,
		Backend:       backend,
		Authenticator: authenticator,
		proxy:         rp,
	}
	return p, nil
}

// NewReverseProxy returns httputil.ReverseProxy object for particular backend
func NewReverseProxy(backend *Backend, transport http.RoundTripper) *httputil.ReverseProxy {
	director := func(r *http.Request) {
		r.URL.Scheme = backend.URL.Scheme
		r.URL.Host = backend.URL.Host
		r.URL.Path = singleJoiningSlash(backend.URL.Path, r.URL.Path)

		if backend.RewriteHost {
			r.Host = backend.URL.Host
		}

		if backend.URL.RawQuery == "" || r.URL.RawQuery == "" {
			r.URL.RawQuery = backend.URL.RawQuery + r.URL.RawQuery
		} else {
			r.URL.RawQuery = backend.URL.RawQuery + "&" + r.URL.RawQuery
		}

		remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		if remoteIP == "" {
			remoteIP = "0.0.0.0"
		}

		r.Header.Set("X-Forwarded-For", remoteIP)
		r.Header.Set("X-Real-IP", remoteIP)
		r.Header.Set("X-Forwarded-Proto", "https")
		r.Header.Set("X-Proxy-App", "svcproxy")
	}

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
}
