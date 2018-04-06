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

func (s *SQLCacheTestSuite) TestMySQLCache() {
	dataSample := []byte("test-byte")
	c, _ := NewSQLCache(s.db)

	err := c.Put(context.Background(), "test-data", dataSample)
	s.Require().NoError(err)

	data, err := c.Get(context.Background(), "test-data")
	s.Require().NoError(err)
	s.Require().NotNil(data)

	err = c.Delete(context.Background(), "test-data")
	s.Require().NoError(err)
}

func (s *SQLCacheTestSuite) SetupTest() {
	s.db = nil
}

func TestMySQLCacheTestSuite(t *testing.T) {
	suite.Run(t, new(SQLCacheTestSuite))
}
