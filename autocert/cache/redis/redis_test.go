package redis

import (
	"context"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/suite"
)

type RedisCacheTestSuite struct {
	suite.Suite
	client *redis.Client
}

func (s *RedisCacheTestSuite) TestMySQLCache() {
	dataSample := []byte("test-byte")
	c, err := NewCache(s.client)
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

func (s *RedisCacheTestSuite) SetupTest() {
	s.client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	s.Require().NotNil(s.client)

	pong, err := s.client.Ping().Result()
	s.Require().NoError(err)
	s.Require().Equal("PONG", pong)
}

func TestRedisCacheTestSuite(t *testing.T) {
	suite.Run(t, new(RedisCacheTestSuite))
}
