package types

import "net/http"

// Middleware interface
type Middleware interface {
	SetOptions(map[string]interface{})
	Middleware(next http.Handler) http.Handler
}
