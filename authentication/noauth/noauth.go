package noauth

import (
	"net/http"

	"github.com/teran/svcproxy/authentication"
)

var _ authentication.Authenticator = &NoAuth{}

// NoAuth implements Authenticator interface
// to provide Basic authenication mechanism(rfc2617)
type NoAuth struct{}

func (nsa *NoAuth) IsAuthenticated(r *http.Request) bool {
	return true
}

func (na *NoAuth) Authenticate(w http.ResponseWriter, r *http.Request) {}
