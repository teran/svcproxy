package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/teran/svcproxy/middleware/filter"
	"github.com/teran/svcproxy/middleware/logging"
	"github.com/teran/svcproxy/middleware/metrics"
	"github.com/teran/svcproxy/middleware/types"
)

var middlewaresMap = map[string]types.Middleware{
	"logging": logging.NewMiddleware(),
	"metrics": metrics.NewMiddleware(),
	"filter":  filter.NewMiddleware(),
}

// Chain allows to chain middlewares dynamically
func Chain(f http.Handler, ms ...map[string]interface{}) http.Handler {
	for _, m := range ms {
		name, ok := m["name"]
		if !ok {
			log.Fatalf("Missed name field in middleware map: %+v", m)
		}

		fm, ok := middlewaresMap[name.(string)]
		if !ok {
			log.Fatalf("Middleware '%s' is requested but not registered.", name.(string))
		}
		log.WithFields(log.Fields{
			"middleware": name,
		}).Debugf("Middleware initialized")

		fm.SetOptions(m)
		f = fm.Middleware(f)
	}

	return f
}
