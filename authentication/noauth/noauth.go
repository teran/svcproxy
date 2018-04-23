package noauth

import (
	"net/http"

	"github.com/teran/svcproxy/authentication"
)

var _ authentication.Authenticator = &NoAuth{}

// NoAuth implements Authenticator interface
// to provide Basic authenication mechanism(rfc2617)
type NoAuth struct{}

// IsAuthenticated for NoAuth Authenticator always returns true.
// It's just a stub to keep the pipeline for requests
func (na *NoAuth) IsAuthenticated(r *http.Request) bool {
	return true
}

// Authenticate for NoAuth Authenticator does nothing.
func (na *NoAuth) Authenticate(w http.ResponseWriter, r *http.Request) {}
