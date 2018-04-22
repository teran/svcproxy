package cache

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"strconv"
	"sync"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/pbkdf2"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"

	// PostgreSQL driver
	_ "github.com/lib/pq"

	sqlcache "svcproxy/autocert/cache/sql"
)

var _ autocert.Cache = &Cache{}

type Cache struct {
	backend       autocert.Cache
	encryptionKey []byte
	usePrecaching bool
	precache      sync.Map
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	if c.usePrecaching {
		data, ok := c.precache.Load(key)
		if ok {
			return data.([]byte), nil
		}
	}

	data, err := c.backend.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if c.encryptionKey == nil {
		c.precache.Store(key, data)
		return data, nil
	}

	data, err = c.decrypt(data)
	if err != nil {
		return nil, err
	}

	c.precache.Store(key, data)

	return data, nil
}

func (c *Cache) Put(ctx context.Context, key string, data []byte) error {
	var resultData []byte
	if c.encryptionKey != nil {
		var err error
		resultData, err = c.encrypt(data)
		if err != nil {
			return err
		}
	} else {
		resultData = data
	}

	err := c.backend.Put(ctx, key, resultData)
	if err != nil {
		return err
	}

	if c.usePrecaching {
		c.precache.Store(key, data)
	}

	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	if c.usePrecaching {
		c.precache.Delete(key)
	}
	return c.backend.Delete(ctx, key)
}

func (c *Cache) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("Ciphertext is too short. Probably corrupted data")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, err
}

func (c *Cache) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// NewCacheFactory returns Cache instance
func NewCacheFactory(backend string, options map[string]string) (autocert.Cache, error) {
	var err error

	var usePrecaching bool
	usePrecachingString, ok := options["usePrecaching"]
	if ok {
		usePrecaching, _ = strconv.ParseBool(usePrecachingString)
	}

	var encKey []byte
	encryptionKey, ok := options["encryptionKey"]
	if ok && encryptionKey != "" {
		h := sha256.New()
		h.Write([]byte(encryptionKey))
		key := h.Sum(nil)
		encKey = pbkdf2.Key(key[:15], key[16:32], 1048, 32, sha256.New)
	}

	var b autocert.Cache
	switch backend {
	case "sql":
		b, err = newSQLCacheBackend(options)
		if err != nil {
			return nil, err
		}
	}

	var c autocert.Cache = &Cache{
		encryptionKey: encKey,
		usePrecaching: usePrecaching,
		precache:      sync.Map{},
		backend:       b,
	}

	return c, err
}

func newSQLCacheBackend(options map[string]string) (autocert.Cache, error) {
	driver, ok := options["driver"]
	if !ok {
		return nil, fmt.Errorf("No driver specified")
	}
	dsn, ok := options["dsn"]
	if !ok {
		return nil, fmt.Errorf("dsn option to backend is required")
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("Error establishing database connection: %s", err)
	}
	if e := db.Ping(); e != nil {
		return nil, fmt.Errorf("Error contacting database: %s", e)
	}
	return sqlcache.NewCache(db)
}
