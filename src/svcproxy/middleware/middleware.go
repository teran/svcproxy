package middleware

import (
	"net/http"

	"svcproxy/middleware/logging"
)

var middlewaresMap = map[string]func(http.Handler) http.Handler{
	"logging": logging.Middleware,
}

// Chain allows to chain middlewares dynamically
func Chain(f http.Handler, ms ...string) http.Handler {
	for _, m := range ms {
		f = middlewaresMap[m](f)
	}

	return f
}
