package sql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

type CacheTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *CacheTestSuite) TestMyCache() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.db, []byte("testKey"))
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
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/svcproxy")
	s.Require().NoError(err)
	s.Require().NotNil(db)
	s.db = db
}

func TestMyCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
