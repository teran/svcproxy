package cache

import (
	"context"
	"testing"
)

func TestMySQLCache(t *testing.T) {
	dataSample := []byte("test-byte")
	c := MySQLCache{}

	if err := c.Put(context.Background(), "test-data", dataSample); err != nil {
		t.Errorf("Error puting data to cache: %s", err)
	}

	if _, err := c.Get(context.Background(), "test-data"); err != nil {
		t.Errorf("Error obtaining data from cache: %s", err)
	}

	if err := c.Delete(context.Background(), "test-data"); err != nil {
		t.Errorf("Error deleting data from cache: %s", err)
	}
}
