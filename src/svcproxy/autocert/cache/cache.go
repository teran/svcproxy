package cache

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/acme/autocert"

	sqlcache "svcproxy/autocert/cache/sql"
)

// NewCacheFactory returns Cache instance
func NewCacheFactory(options map[string]string) (autocert.Cache, error) {
	driver, ok := options["driver"]
	if !ok {
		return nil, fmt.Errorf("No driver specified")
	}

	db, err := sql.Open(driver, options["dsn"])
	if err != nil {
		log.Fatalf("Error establising database connection: %s", err)
	}

	return sqlcache.NewCache(db, []byte(options["encryptionKey"]))
}
