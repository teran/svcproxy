package cache

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/acme/autocert"
)

var _ autocert.Cache = &SQLCache{}

// SQLCache implements autocert.Cache with MySQL database
type SQLCache struct {
	db *sql.DB
}

// NewSQLCache returns SQLCache instance
func NewSQLCache(db *sql.DB) (*SQLCache, error) {
	return &SQLCache{
		db: db,
	}, nil
}

// Get retrieves certificate data from cache
func (m *SQLCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

// Put stores certificate data to cache
func (m *SQLCache) Put(ctx context.Context, key string, data []byte) error {
	return nil
}

// Delete removes certificate data from cache
func (m *SQLCache) Delete(ctx context.Context, key string) error {
	return nil
}
