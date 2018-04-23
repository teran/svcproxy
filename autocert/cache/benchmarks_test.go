package cache

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

// SQL/MySQL benchmarks
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

func BenchmarkCacheGetSQLMySQLWithEncryptionAndWithoutPrecaching(b *testing.B) {
	r := require.New(b)
	options := map[string]string{
		"driver":        "mysql",
		"dsn":           "root@tcp(127.0.0.1:3306)/svcproxy?parseTime=true",
		"usePrecaching": "false",
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

// SQL/PostgreSQL benchmarks
func BenchmarkCacheGetSQLPostgreSQLWithoutEncryptionAndPrecaching(b *testing.B) {
	r := require.New(b)
	options := map[string]string{
		"driver":        "postgres",
		"dsn":           "postgres://postgres@localhost/svcproxy?sslmode=disable",
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

func BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndWithoutPrecaching(b *testing.B) {
	r := require.New(b)
	options := map[string]string{
		"driver":        "postgres",
		"dsn":           "postgres://postgres@localhost/svcproxy?sslmode=disable",
		"usePrecaching": "false",
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

func BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndPrecaching(b *testing.B) {
	r := require.New(b)
	options := map[string]string{
		"driver":        "postgres",
		"dsn":           "postgres://postgres@localhost/svcproxy?sslmode=disable",
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

// Dir benchmarks
func BenchmarkCacheGetDirWithoutEncryptionAndPrecaching(b *testing.B) {
	r := require.New(b)

	dir, err := ioutil.TempDir("", "dircache_test")
	r.NoError(err)

	options := map[string]string{
		"path":          dir,
		"usePrecaching": "false",
		"encryptionKey": "",
	}
	c, err := NewCacheFactory("dir", options)
	r.NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	r.NoError(err)

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "test-key")
	}
}

func BenchmarkCacheGetDirWithEncryptionAndWithoutPrecaching(b *testing.B) {
	r := require.New(b)

	dir, err := ioutil.TempDir("", "dircache_test")
	r.NoError(err)

	options := map[string]string{
		"path":          dir,
		"usePrecaching": "false",
		"encryptionKey": "blah",
	}
	c, err := NewCacheFactory("dir", options)
	r.NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	r.NoError(err)

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "test-key")
	}
}

func BenchmarkCacheGetDirWithEncryptionAndPrecaching(b *testing.B) {
	r := require.New(b)

	dir, err := ioutil.TempDir("", "dircache_test")
	r.NoError(err)

	options := map[string]string{
		"path":          dir,
		"usePrecaching": "true",
		"encryptionKey": "blah",
	}
	c, err := NewCacheFactory("dir", options)
	r.NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	r.NoError(err)

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "test-key")
	}
}

// Redis benchmarks
func BenchmarkCacheGetRedisWithoutEncryptionAndPrecaching(b *testing.B) {
	r := require.New(b)

	options := map[string]string{
		"addr":          "127.0.0.1:6379",
		"usePrecaching": "false",
		"encryptionKey": "",
	}
	c, err := NewCacheFactory("redis", options)
	r.NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	r.NoError(err)

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "test-key")
	}
}

func BenchmarkCacheGetRedisWithEncryptionAndWithoutPrecaching(b *testing.B) {
	r := require.New(b)

	options := map[string]string{
		"addr":          "127.0.0.1:6379",
		"usePrecaching": "false",
		"encryptionKey": "blah",
	}
	c, err := NewCacheFactory("redis", options)
	r.NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	r.NoError(err)

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "test-key")
	}
}

func BenchmarkCacheGetRedisWithEncryptionAndPrecaching(b *testing.B) {
	r := require.New(b)

	options := map[string]string{
		"addr":          "127.0.0.1:6379",
		"usePrecaching": "true",
		"encryptionKey": "blah",
	}
	c, err := NewCacheFactory("redis", options)
	r.NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	r.NoError(err)

	for n := 0; n < b.N; n++ {
		c.Get(context.Background(), "test-key")
	}
}
