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

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/pbkdf2"

	"svcproxy/autocert/cache/sql/mysql"
	"svcproxy/autocert/cache/sql/postgresql"
)

var _ autocert.Cache = &Cache{}

type Driver interface {
	Get(key string) ([]byte, error)
	Put(key string, data []byte) error
	Delete(key string) error
}

// Cache implements autocert.Cache with MySQL database
type Cache struct {
	driver        Driver
	encryptionKey []byte
}

// NewCache returns Cache instance
func NewCache(db *sql.DB, encryptionKey []byte) (*Cache, error) {
	h := sha256.New()
	h.Write(encryptionKey)
	key := h.Sum(nil)

	var driver Driver

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

	return &Cache{
		driver:        driver,
		encryptionKey: pbkdf2.Key(key[:15], key[16:32], 1048, 32, sha256.New),
	}, nil
}

// Get retrieves certificate data from cache
func (m *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := m.driver.Get(key)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: add test for ErrCacheMiss
			return nil, autocert.ErrCacheMiss
		}
		return nil, err
	}

	decryptedData, err := m.decrypt(data)
	if err != nil {
		return nil, err
	}

	return decryptedData, err
}

// Put stores certificate data to cache
func (m *Cache) Put(ctx context.Context, key string, data []byte) error {
	encryptedData, err := m.encrypt(data)
	if err != nil {
		return err
	}

	return m.driver.Put(key, encryptedData)
}

// Delete removes certificate data from cache
func (m *Cache) Delete(ctx context.Context, key string) error {
	return m.driver.Delete(key)
}

func (m *Cache) decrypt(ciphertext []byte) ([]byte, error) {
	ct := make([]byte, base64.StdEncoding.DecodedLen(len(ciphertext)))
	l, err := base64.StdEncoding.Decode(ct, ciphertext)
	if err != nil {
		return nil, err
	}

	ciphertext = ct[:l]

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

	ct := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(ct, ciphertext)
	if err != nil {
		return nil, err
	}

	return ct, nil
}
