package middleware

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/teran/svcproxy/middleware/filter"
	"github.com/teran/svcproxy/middleware/gzip"
	"github.com/teran/svcproxy/middleware/logging"
	"github.com/teran/svcproxy/middleware/metrics"
	"github.com/teran/svcproxy/middleware/types"
)

type middlewareDefinition struct {
	middleware types.Middleware
	config     types.MiddlewareConfig
}

var middlewaresMap = map[string]middlewareDefinition{
	"filter": middlewareDefinition{
		middleware: filter.NewMiddleware(),
		config:     &filter.Config{},
	},
	"gzip": middlewareDefinition{
		middleware: gzip.NewMiddleware(),
		config:     &gzip.GzipConfig{},
	},
	"logging": middlewareDefinition{
		middleware: logging.NewMiddleware(),
	},
	"metrics": middlewareDefinition{
		middleware: metrics.NewMiddleware(),
	},
}

// Chain allows to chain middlewares dynamically
func Chain(f http.Handler, ms ...map[string]interface{}) (http.Handler, error) {
	for _, m := range ms {
		name, ok := m["name"]
		if !ok {
			return nil, fmt.Errorf("Missed name field in middleware map: %+v", m)
		}

		md, ok := middlewaresMap[name.(string)]
		if !ok {
			return nil, fmt.Errorf("middleware `%s` is requested but not registered", name.(string))
		}
		log.WithFields(log.Fields{
			"middleware": name,
		}).Debugf("Middleware initialized")

		if md.config != nil {
			err := md.config.Unpack(m)
			if err != nil {
				return nil, err
			}

			err = md.middleware.SetConfig(md.config)
			if err != nil {
				return nil, err
			}
		}
		f = md.middleware.Middleware(f)
	}

	return f, nil
}
