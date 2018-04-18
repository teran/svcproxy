package cache

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
}

func (s *CacheTestSuite) TestInitializeSQLCache() {
	options := map[string]string{
		"driver":        "mysql",
		"dsn":           "root@tcp(127.0.0.1:3306)/svcproxy?parseTime=true",
		"encryptionKey": "testkey",
	}
	_, err := NewCacheFactory("sql", options)
	s.Require().NoError(err)
}

func (s *CacheTestSuite) TestInitializeSQLCacheNoEncryption() {
	options := map[string]string{
		"driver": "mysql",
		"dsn":    "root@tcp(127.0.0.1:3306)/svcproxy?parseTime=true",
	}
	_, err := NewCacheFactory("sql", options)
	s.Require().NoError(err)
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
