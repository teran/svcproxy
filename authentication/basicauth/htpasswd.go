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
	passwdFile string
}

// IsValidCredentials checks client's credentials against htpasswd file
func (h *HTPasswd) IsValidCredentials(username, password string) (bool, error) {
	fp, err := os.Open(h.passwdFile)
	if err != nil {
		return false, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		credentials := strings.Split(line, ":")
		if len(credentials) != 2 {
			return false, fmt.Errorf("BasicAuth htpasswd authenticator error: passwd file '%s' mailformed", h.passwdFile)
		}
		if username == credentials[0] {
			if err := bcrypt.CompareHashAndPassword([]byte(credentials[1]), []byte(password)); err == nil {
				return true, nil
			}
			return false, nil
		}
	}

	return false, nil
}
