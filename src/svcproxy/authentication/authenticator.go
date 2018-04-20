package authentication

import "net/http"

// Authenticator is used by svcproxy to authenticate requests to services
type Authenticator interface {
	IsAuthenticated(r *http.Request) bool
	Authenticate(w http.ResponseWriter, r *http.Request)
}
