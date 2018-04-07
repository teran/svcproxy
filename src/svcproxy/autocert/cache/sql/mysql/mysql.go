package mysql

import (
	"database/sql"
)

// MySQL database driver abstraction
type MySQL struct {
	DB *sql.DB
}

// Get serves to retrieve cached data from MySQL database
func (m *MySQL) Get(key string) ([]byte, error) {
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
		return nil, err
	}

	return value, nil
}

// Put serves to place data to MySQL database as cache
func (m *MySQL) Put(key string, data []byte) error {
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
func (m *MySQL) Delete(key string) error {
	_, err := m.DB.Exec(`
		DELETE FROM
			autocert_cache
		WHERE
			cache_key = ?
	`, key)

	return err
}
