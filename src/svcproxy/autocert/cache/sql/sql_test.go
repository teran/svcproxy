package sql

import (
	"context"
	"crypto/rand"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"

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

func (s *SQLCacheTestSuite) TestMySQLCacheNoEncryption() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.mysql, nil, false)
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

func (s *SQLCacheTestSuite) TestMySQLCache() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.mysql, []byte("testKey"), false)
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

func (s *SQLCacheTestSuite) TestPostgreSQLCache() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.postgresql, []byte("testKey"), false)
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

func (s *SQLCacheTestSuite) SetupTest() {
	var err error
	s.mysql, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/svcproxy")
	s.Require().NoError(err)
	s.Require().NotNil(s.mysql)

	s.postgresql, err = sql.Open("postgres", "postgres://postgres@localhost/svcproxy?sslmode=disable")
	s.Require().NoError(err)
	s.Require().NotNil(s.postgresql)
}

func (s *SQLCacheTestSuite) TestMySQLCacheWithPrecaching() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.mysql, []byte("testKey"), true)
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

func TestSQLCacheTestSuite(t *testing.T) {
	suite.Run(t, new(SQLCacheTestSuite))
}

func BenchmarkGetFromCacheMySQL(b *testing.B) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/svcproxy")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	c, err := NewCache(db, []byte("testKey"), false)
	if err != nil {
		panic(err)
	}

	dataSample := make([]byte, 4096)
	rand.Read(dataSample)

	err = c.Put(context.Background(), "testdata", dataSample)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "testdata")
	}
}

func BenchmarkGetFromCacheMySQLNoEncryption(b *testing.B) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/svcproxy")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	c, err := NewCache(db, nil, false)
	if err != nil {
		panic(err)
	}

	dataSample := make([]byte, 4096)
	rand.Read(dataSample)

	err = c.Put(context.Background(), "testdata", dataSample)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "testdata")
	}
}

func BenchmarkGetFromCachePostgreSQL(b *testing.B) {
	db, err := sql.Open("postgres", "postgres://postgres@localhost/svcproxy?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}
	c, err := NewCache(db, []byte("testKey"), false)
	if err != nil {
		panic(err)
	}

	dataSample := make([]byte, 4096)
	rand.Read(dataSample)

	err = c.Put(context.Background(), "test-data", dataSample)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "test-data")
	}
}

func BenchmarkGetFromCachePostgreSQLNoEncryption(b *testing.B) {
	db, err := sql.Open("postgres", "postgres://postgres@localhost/svcproxy?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}
	c, err := NewCache(db, nil, false)
	if err != nil {
		panic(err)
	}

	dataSample := make([]byte, 4096)
	rand.Read(dataSample)

	err = c.Put(context.Background(), "testdata_unencrypted", dataSample)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "testdata_unencrypted")
	}
}

func BenchmarkGetFromCacheMySQLWithPrecaching(b *testing.B) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/svcproxy")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	c, err := NewCache(db, []byte("testKey"), true)
	if err != nil {
		panic(err)
	}

	dataSample := make([]byte, 4096)
	rand.Read(dataSample)

	err = c.Put(context.Background(), "testdata", dataSample)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "testdata")
	}
}
