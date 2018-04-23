package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"

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
		log.WithFields(log.Fields{
			"middleware": name,
		}).Debugf("Middleware initialized")
		f = fm(f, m)
	}

	return f
}
