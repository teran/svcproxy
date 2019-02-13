package types

import "net/http"

// Middleware interface
type Middleware interface {
	SetConfig(MiddlewareConfig) error
	Middleware(next http.Handler) http.Handler
}

// MiddlewareConfig interface
type MiddlewareConfig interface {
	Unpack(map[string]interface{}) error
}
