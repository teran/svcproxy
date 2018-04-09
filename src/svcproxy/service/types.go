package service

import (
	"net/http/httputil"
	"net/url"
)

// Service interface
type Service interface {
	AddProxy(*Proxy) error
}

// Proxy type
type Proxy struct {
	Frontend *Frontend
	Backend  *Backend
	proxy    *httputil.ReverseProxy
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
