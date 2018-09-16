package sql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/acme/autocert"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

type SQLCacheTestSuite struct {
	suite.Suite
	mysql      *sql.DB
	postgresql *sql.DB
}

func (s *SQLCacheTestSuite) TestMySQLCache() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.mysql)
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

	data, err = c.Get(context.Background(), "test-data")
	s.Require().Error(err)
	s.Require().Equal(autocert.ErrCacheMiss, err)
	s.Require().Nil(data)
}

func (s *SQLCacheTestSuite) TestPostgreSQLCache() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.postgresql)
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

	data, err = c.Get(context.Background(), "test-data")
	s.Require().Error(err)
	s.Require().Equal(autocert.ErrCacheMiss, err)
	s.Require().Nil(data)
}

func (s *SQLCacheTestSuite) SetupTest() {
	var err error
	s.mysql, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/svcproxy?parseTime=true")
	s.Require().NoError(err)
	s.Require().NotNil(s.mysql)

	_, err = s.mysql.Exec("DROP TABLE IF EXISTS `svcproxy`;")
	s.Require().NoError(err)

	// Migrate manually to avoid possible race on each NewCache call
	err = maybeMigrate(s.mysql, "mysql")
	s.Require().NoError(err)

	s.postgresql, err = sql.Open("postgres", "postgres://postgres@localhost/svcproxy?sslmode=disable")
	s.Require().NoError(err)
	s.Require().NotNil(s.postgresql)

	_, err = s.postgresql.Exec("DROP TABLE IF EXISTS svcproxy;")
	s.Require().NoError(err)

	// Migrate manually to avoid possible race on each NewCache call
	err = maybeMigrate(s.postgresql, "postgres")
	s.Require().NoError(err)
}

func TestSQLCacheTestSuite(t *testing.T) {
	suite.Run(t, new(SQLCacheTestSuite))
}
