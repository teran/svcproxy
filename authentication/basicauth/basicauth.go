package basicauth

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/teran/svcproxy/authentication"
)

var _ authentication.Authenticator = &BasicAuth{}

// Backend is an interface for underlying authentication storages
// example for such storage could be: htpasswd file, sql database, PAM, etc.
// The main purpose for types implementing BasicAuthBackend is to validate
// credentials provided by user against specific storage.
type Backend interface {
	IsValidCredentials(username, password string) (bool, error)
}

// BasicAuth implements Authenticator interface
// to provide Basic authenication mechanism(rfc2617)
type BasicAuth struct {
	passwdFile string
	backend    Backend
}

// NewBasicAuth creates new BasicAuth object with requested backend
func NewBasicAuth(backendName string, options map[string]string) (authentication.Authenticator, error) {
	backend, err := NewBasicAuthBackend(backendName, options)
	if err != nil {
		return nil, err
	}

	return &BasicAuth{
		backend: backend,
	}, nil
}

// NewBasicAuthBackend creates new Backend by name
func NewBasicAuthBackend(name string, options map[string]string) (Backend, error) {
	switch name {
	case "htpasswd":
		passwdFile, ok := options["file"]
		if !ok {
			return nil, fmt.Errorf("'file' option must be passed to htpasswd backend but not specified")
		}
		return NewHTPasswdBackend(passwdFile)
	}
	return nil, fmt.Errorf("Unknown backend: %s", name)
}

// IsAuthenticated verifies credentials passed in request(if any) against Backend
func (ba *BasicAuth) IsAuthenticated(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}

	ok, err := ba.backend.IsValidCredentials(username, password)
	if err != nil {
		log.Warnf("Error verifying credentials: %s", err)
		return false
	}
	if !ok {
		return false
	}
	return true
}

// Authenticate in BasicAuth authenticator simply sends headers to client
// to forse them to show HTTP Basic Auth login form.
func (ba *BasicAuth) Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted area"`)
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
