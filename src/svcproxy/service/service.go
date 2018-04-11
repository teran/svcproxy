package service

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
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
			redir, err := url.Parse(fmt.Sprintf("https://%s%s", host, r.URL.Path))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			http.Redirect(w, r, redir.String(), http.StatusFound)
			return
		}
	}

	p.proxy.ServeHTTP(w, r)
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
