package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkCacheGetSQLMySQLWithoutEncryptionAndPrecaching(b *testing.B) {
	r := require.New(b)
	options := map[string]string{
		"driver":        "mysql",
		"dsn":           "root@tcp(127.0.0.1:3306)/svcproxy?parseTime=true",
		"usePrecaching": "false",
		"encryptionKey": "",
	}
	c, err := NewCacheFactory("sql", options)
	r.NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	r.NoError(err)

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "test-key")
	}
}

func BenchmarkCacheGetSQLMySQLWithEncryptionAndPrecaching(b *testing.B) {
	r := require.New(b)
	options := map[string]string{
		"driver":        "mysql",
		"dsn":           "root@tcp(127.0.0.1:3306)/svcproxy?parseTime=true",
		"usePrecaching": "true",
		"encryptionKey": "blah",
	}
	c, err := NewCacheFactory("sql", options)
	r.NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	r.NoError(err)

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "test-key")
	}
}
