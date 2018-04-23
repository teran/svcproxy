package middleware

import (
	"log"
	"net/http"

	"github.com/teran/svcproxy/middleware/logging"
	"github.com/teran/svcproxy/middleware/metrics"
)

var middlewaresMap = map[string]func(http.Handler, map[string]string) http.Handler{
	"logging": logging.Middleware,
	"metrics": metrics.NewMetricsMiddleware().Middleware,
}

// Chain allows to chain middlewares dynamically
func Chain(f http.Handler, ms ...map[string]string) http.Handler {
	for _, m := range ms {
		name, ok := m["name"]
		if !ok {
			log.Fatalf("Missed name field in middleware map: %+v", m)
		}
		fm, ok := middlewaresMap[name]
		if !ok {
			log.Fatalf("Middleware '%s' is requested but not registered.", name)
		}
		log.Printf("Chaining middleware %s", name)
		f = fm(f, m)
	}

	return f
}
