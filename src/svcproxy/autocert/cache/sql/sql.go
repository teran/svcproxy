package sql

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"sync"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/pbkdf2"

	"svcproxy/autocert/cache/sql/mysql"
	"svcproxy/autocert/cache/sql/postgresql"
)

var _ autocert.Cache = &Cache{}

// Cache implements autocert.Cache with MySQL database
type Cache struct {
	driver        autocert.Cache
	encryptionKey []byte
	usePrecaching bool
	precache      sync.Map
}

// NewCache returns Cache instance
func NewCache(db *sql.DB, encryptionKey []byte, usePrecaching bool) (*Cache, error) {
	h := sha256.New()
	h.Write(encryptionKey)
	key := h.Sum(nil)

	var driver autocert.Cache

	switch fmt.Sprintf("Driver: %s", reflect.TypeOf(db.Driver())) {
	case "Driver: *mysql.MySQLDriver":
		driver = &mysql.MySQL{
			DB: db,
		}
	case "Driver: *pq.Driver":
		driver = &postgresql.PostgreSQL{
			DB: db,
		}
	default:
		return nil, fmt.Errorf("Unsupported driver")
	}

	var encKey []byte
	if encryptionKey != nil {
		encKey = pbkdf2.Key(key[:15], key[16:32], 1048, 32, sha256.New)
	}

	return &Cache{
		driver:        driver,
		encryptionKey: encKey,
		usePrecaching: usePrecaching,
		precache:      sync.Map{},
	}, nil
}

// Get retrieves certificate data from cache
func (m *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	if m.usePrecaching {
		data, ok := m.precache.Load(key)
		if ok {
			return data.([]byte), nil
		}
	}

	data, err := m.driver.Get(ctx, key)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: add test for ErrCacheMiss
			return nil, autocert.ErrCacheMiss
		}
		return nil, err
	}

	data, err = m.decode(data)
	if err != nil {
		return nil, err
	}

	if m.encryptionKey == nil {
		return data, nil
	}

	data, err = m.decrypt(data)
	if err != nil {
		return nil, err
	}

	m.precache.Store(key, data)

	return data, err
}

// Put stores certificate data to cache
func (m *Cache) Put(ctx context.Context, key string, data []byte) error {
	if m.usePrecaching {
		m.precache.Store(key, data)
	}

	if m.encryptionKey != nil {
		var err error
		data, err = m.encrypt(data)
		if err != nil {
			m.precache.Delete(key)
			return err
		}
	}

	data = m.encode(data)
	err := m.driver.Put(ctx, key, data)
	if err != nil {
		m.precache.Delete(key)
		return err
	}

	return nil
}

// Delete removes certificate data from cache
func (m *Cache) Delete(ctx context.Context, key string) error {
	if m.usePrecaching {
		m.precache.Delete(key)
	}

	return m.driver.Delete(ctx, key)
}

func (m *Cache) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.encryptionKey)
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

func (m *Cache) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.encryptionKey)
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

func (m *Cache) decode(input []byte) ([]byte, error) {
	ct := make([]byte, base64.StdEncoding.DecodedLen(len(input)))
	l, err := base64.StdEncoding.Decode(ct, input)
	if err != nil {
		return nil, err
	}

	return ct[:l], nil
}

func (m *Cache) encode(input []byte) []byte {
	ct := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
	base64.StdEncoding.Encode(ct, input)

	return ct
}
