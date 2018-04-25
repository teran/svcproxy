package basicauth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var _ Backend = &HTPasswd{}

// HTPasswd backend type
type HTPasswd struct {
	cache map[string]string
}

// NewHTPasswdBackend returns HTPasswd instance
func NewHTPasswdBackend(passwdFile string) (Backend, error) {
	h := &HTPasswd{
		cache: make(map[string]string),
	}

	fp, err := os.Open(passwdFile)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		credentials := strings.Split(line, ":")
		if len(credentials) != 2 {
			return nil, fmt.Errorf("BasicAuth htpasswd authenticator error: passwd file '%s' mailformed", passwdFile)
		}

		h.cache[credentials[0]] = credentials[1]
	}

	return h, nil
}

// IsValidCredentials checks client's credentials against htpasswd file
func (h *HTPasswd) IsValidCredentials(username, password string) (bool, error) {
	cp, ok := h.cache[username]
	if !ok {
		return false, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cp), []byte(password)); err == nil {
		return true, nil
	}

	return false, nil
}
