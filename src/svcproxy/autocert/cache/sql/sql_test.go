package sql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

type CacheTestSuite struct {
	suite.Suite
	mysql      *sql.DB
	postgresql *sql.DB
}

func (s *CacheTestSuite) TestMySQLCache() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.mysql, []byte("testKey"))
	s.Require().NoError(err)
	s.Require().NotNil(c)

	err = c.Put(context.Background(), "test-data", dataSample)
	s.Require().NoError(err)

	data, err := c.Get(context.Background(), "test-data")
	s.Require().NoError(err)
	s.Require().NotNil(data)
	s.Require().Equal(dataSample, data)

	err = c.Delete(context.Background(), "test-data")
	s.Require().NoError(err)
}

func (s *CacheTestSuite) TestPostgreSQLCache() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.postgresql, []byte("testKey"))
	s.Require().NoError(err)
	s.Require().NotNil(c)

	err = c.Put(context.Background(), "test-data", dataSample)
	s.Require().NoError(err)

	data, err := c.Get(context.Background(), "test-data")
	s.Require().NoError(err)
	s.Require().NotNil(data)
	s.Require().Equal(dataSample, data)

	err = c.Delete(context.Background(), "test-data")
	s.Require().NoError(err)
}

func (s *CacheTestSuite) SetupTest() {
	var err error
	s.mysql, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/svcproxy")
	s.Require().NoError(err)
	s.Require().NotNil(s.mysql)

	s.postgresql, err = sql.Open("postgres", "postgres://postgres@localhost/svcproxy?sslmode=disable")
	s.Require().NoError(err)
	s.Require().NotNil(s.postgresql)
}

func TestMyCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
