package cache

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SQLCacheTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *SQLCacheTestSuite) TestMySQLCache(t *testing.T) {
	dataSample := []byte("test-byte")
	c, _ := NewSQLCache(s.db)

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

func (s *SQLCacheTestSuite) SetupTest() {
	s.db = nil
}

func TestMySQLCacheTestSuite(t *testing.T) {
	suite.Run(t, new(SQLCacheTestSuite))
}
