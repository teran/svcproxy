package factory

import (
	"fmt"

	"github.com/teran/svcproxy/authentication"
	ba "github.com/teran/svcproxy/authentication/basicauth"
	na "github.com/teran/svcproxy/authentication/noauth"
)

// NewAuthenticator returns specific authenticator based on configuration
func NewAuthenticator(name string, options map[string]string) (authentication.Authenticator, error) {
	switch name {
	// No authentication by default
	case "":
		return &na.NoAuth{}, nil
	// NoAuth is the same
	case "NoAuth":
		return &na.NoAuth{}, nil
	case "BasicAuth":
		backend, ok := options["backend"]
		if !ok {
			return nil, fmt.Errorf("backend option is required for BasicAuth authenticator")
		}
		return ba.NewBasicAuth(backend, options)
	}
	return nil, fmt.Errorf("Unknown authenticator %s", name)
}
