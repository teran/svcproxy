package service

import (
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
}

// Frontend type
type Frontend struct {
	FQDN string
}

// Backend type
type Backend struct {
	URL         *url.URL
	RewriteHost bool
}
