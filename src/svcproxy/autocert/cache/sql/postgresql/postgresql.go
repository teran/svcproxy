package postgresql

import (
	"database/sql"
)

// PostgreSQL database driver abstraction
type PostgreSQL struct {
	db *sql.DB
}

func (m *PostgreSQL) Get(key string) ([]byte, error) {
	var value []byte

	err := m.db.QueryRow(`
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

func (m *PostgreSQL) Put(key string, data []byte) error {
	_, err := m.db.Exec(`
		INSERT INTO
			autocert_cache
			(cache_key, cache_value)
		VALUES
			($1, $2)
		ON
			DUPLICATE KEY
		UPDATE
			cache_key=VALUES(cache_key),
			cache_value=VALUES(cache_value)
	`, key, data)

	return err
}

func (m *PostgreSQL) Delete(key string) error {
	_, err := m.db.Exec(`
		DELETE FROM
			autocert_cache
		WHERE
			cache_key = $1
	`, key)

	return err
}
