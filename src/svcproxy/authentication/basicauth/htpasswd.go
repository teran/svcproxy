package basicauth

import (
	"bufio"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var _ BasicAuthBackend = &HTPasswd{}

type HTPasswd struct {
	passwdFile string
}

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
			log.Printf("BasicAuth authenticator error: passwd file %s mailformed", h.passwdFile)
			continue
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
