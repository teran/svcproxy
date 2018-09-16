package mysql

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/acme/autocert"
)

var _ autocert.Cache = &MySQL{}

// MySQL database driver abstraction
type MySQL struct {
	DB *sql.DB
}

// Get serves to retrieve cached data from MySQL database
func (m *MySQL) Get(ctx context.Context, key string) ([]byte, error) {
	var value []byte

	err := m.DB.QueryRow(`
		SELECT
			cache_value
		FROM
			autocert_cache
		WHERE
			cache_key = ?
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

// Put serves to place data to MySQL database as cache
func (m *MySQL) Put(ctx context.Context, key string, data []byte) error {
	_, err := m.DB.Exec(`
		INSERT INTO
			autocert_cache
			(cache_key, cache_value)
		VALUES
			(?, ?)
		ON
			DUPLICATE KEY
		UPDATE
			cache_key=VALUES(cache_key),
			cache_value=VALUES(cache_value)
	`, key, data)

	return err
}

// Delete serves to delete data from MySQL database
func (m *MySQL) Delete(ctx context.Context, key string) error {
	_, err := m.DB.Exec(`
		DELETE FROM
			autocert_cache
		WHERE
			cache_key = ?
	`, key)

	return err
}
