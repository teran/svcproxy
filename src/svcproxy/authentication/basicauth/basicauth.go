package basicauth

import (
	"fmt"
	"log"
	"net/http"
	"svcproxy/authentication"
)

var _ authentication.Authenticator = &BasicAuth{}

type BasicAuthBackend interface {
	IsValidCredentials(username, password string) (bool, error)
}

// BasicAuth implements Authenticator interface
// to provide Basic authenication mechanism(rfc2617)
type BasicAuth struct {
	passwdFile string
	backend    BasicAuthBackend
}

func NewBasicAuth(backendName string, options map[string]string) (authentication.Authenticator, error) {
	backend, err := NewBasicAuthBackend(backendName, options)
	if err != nil {
		return nil, err
	}

	return &BasicAuth{
		backend: backend,
	}, nil
}

func NewBasicAuthBackend(name string, options map[string]string) (BasicAuthBackend, error) {
	switch name {
	case "htpasswd":
		passwdFile, ok := options["file"]
		if !ok {
			return nil, fmt.Errorf("file option must be passed to htpasswd backend but not specified")
		}
		return &HTPasswd{
			passwdFile: passwdFile,
		}, nil
	}
	return nil, fmt.Errorf("Unknown backend: %s", name)
}

func (ba *BasicAuth) IsAuthenticated(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}

	ok, err := ba.backend.IsValidCredentials(username, password)
	if err != nil {
		log.Printf("Error verifying credentials: %s", err)
		return false
	}
	if !ok {
		return false
	}
	return true
}

func (ba *BasicAuth) Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted area"`)
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
