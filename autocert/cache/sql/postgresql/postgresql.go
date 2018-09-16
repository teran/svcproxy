package postgresql

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/acme/autocert"
)

var _ autocert.Cache = &PostgreSQL{}

// PostgreSQL database driver abstraction
type PostgreSQL struct {
	DB *sql.DB
}

// Get serves to retrieve cached data from MySQL database
func (m *PostgreSQL) Get(ctx context.Context, key string) ([]byte, error) {
	var value []byte

	err := m.DB.QueryRow(`
		SELECT
			cache_value
		FROM
			autocert_cache
		WHERE
			cache_key = $1
		LIMIT 1
	`, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, autocert.ErrCacheMiss
		}
		return nil, err
	}

	return value, nil
}

// Put serves to place data to PostgreSQL database as cache
func (m *PostgreSQL) Put(ctx context.Context, key string, data []byte) error {
	_, err := m.DB.Exec(`
		INSERT INTO
			autocert_cache
			(cache_key, cache_value)
		VALUES
			($1, $2)
		ON CONFLICT (cache_key) DO UPDATE
		SET
			cache_key = $1,
			cache_value = $2
	`, key, data)

	return err
}

// Delete serves to delete data from MySQL database
func (m *PostgreSQL) Delete(ctx context.Context, key string) error {
	_, err := m.DB.Exec(`
		DELETE FROM
			autocert_cache
		WHERE
			cache_key = $1
	`, key)

	return err
}
