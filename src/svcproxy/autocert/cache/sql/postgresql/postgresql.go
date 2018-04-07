package postgresql

import (
	"database/sql"
)

// PostgreSQL database driver abstraction
type PostgreSQL struct {
	DB *sql.DB
}

// Get serves to retrieve cached data from MySQL database
func (m *PostgreSQL) Get(key string) ([]byte, error) {
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
		return nil, err
	}

	return value, nil
}

// Put serves to place data to PostgreSQL database as cache
func (m *PostgreSQL) Put(key string, data []byte) error {
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
func (m *PostgreSQL) Delete(key string) error {
	_, err := m.DB.Exec(`
		DELETE FROM
			autocert_cache
		WHERE
			cache_key = $1
	`, key)

	return err
}
