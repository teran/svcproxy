package cache

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

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/pbkdf2"
)

var _ autocert.Cache = &SQLCache{}

// SQLCache implements autocert.Cache with MySQL database
type SQLCache struct {
	db            *sql.DB
	encryptionKey []byte
}

// NewSQLCache returns SQLCache instance
func NewSQLCache(db *sql.DB, encryptionKey []byte) (*SQLCache, error) {
	h := sha256.New()
	h.Write(encryptionKey)
	key := h.Sum(nil)

	return &SQLCache{
		db:            db,
		encryptionKey: pbkdf2.Key(key[:15], key[16:32], 1048, 32, sha256.New),
	}, nil
}

// Get retrieves certificate data from cache
func (m *SQLCache) Get(ctx context.Context, key string) ([]byte, error) {
	var value []byte

	err := m.db.QueryRow(`
		SELECT
			cache_value
		FROM
			autocert_cache
		WHERE
			cache_key = ?
		LIMIT 1
	`, key).Scan(&value)
	if err != nil {
		return nil, err
	}

	decryptedData, err := m.decrypt(value)
	if err != nil {
		return nil, err
	}

	return decryptedData, err
}

// Put stores certificate data to cache
func (m *SQLCache) Put(ctx context.Context, key string, data []byte) error {
	encryptedData, err := m.encrypt(data)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`
		INSERT INTO
			autocert_cache
			(cache_key, cache_value)
		VALUES
			(?, ?)
		ON
			DUPLICATE KEY
		UPDATE
			cache_key=VALUES(cache_key),
			cache_value=VALUES(cache_value)
	`, key, encryptedData)

	return err
}

// Delete removes certificate data from cache
func (m *SQLCache) Delete(ctx context.Context, key string) error {
	_, err := m.db.Exec(`
		DELETE FROM
			autocert_cache
		WHERE
			cache_key = ?
	`, key)
	return err
}

func (m *SQLCache) decrypt(ciphertext []byte) ([]byte, error) {
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

func (m *SQLCache) encrypt(plaintext []byte) ([]byte, error) {
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
