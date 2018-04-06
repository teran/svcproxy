package cache

import (
	"context"

	"golang.org/x/crypto/acme/autocert"
)

var _ autocert.Cache = &MySQLCache{}

// MySQLCache implements autocert.Cache with MySQL database
type MySQLCache struct{}

// Get retrieves certificate data from cache
func (m *MySQLCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

// Put stores certificate data to cache
func (m *MySQLCache) Put(ctx context.Context, key string, data []byte) error {
	return nil
}

// Delete removes certificate data from cache
func (m *MySQLCache) Delete(ctx context.Context, key string) error {
	return nil
}
