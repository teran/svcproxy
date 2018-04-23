package service

import (
	"net/http/httputil"
	"net/url"

	"github.com/teran/svcproxy/authentication"
)

// Service interface
type Service interface {
	AddProxy(*Proxy) error
}

// Proxy type
type Proxy struct {
	Frontend      *Frontend
	Backend       *Backend
	proxy         *httputil.ReverseProxy
	Authenticator authentication.Authenticator
}

// Frontend type
type Frontend struct {
	FQDN                string
	HTTPHandler         string
	ResponseHTTPHeaders map[string]string
}

// Backend type
type Backend struct {
	URL         *url.URL
	RewriteHost bool
}
