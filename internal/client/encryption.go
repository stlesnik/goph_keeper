package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

// ClientEncryption handles encryption on client side
type ClientEncryption struct {
	userKey []byte
}

// NewClientEncryption creates new client encryption instance
func NewClientEncryption(password string) *ClientEncryption {
	salt := []byte("goph_keeper_salt_2024")
	userKey := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	return &ClientEncryption{
		userKey: userKey,
	}
}

// EncryptData encrypts data on client side
func (ce *ClientEncryption) EncryptData(data string) (string, string, error) {
	block, err := aes.NewCipher(ce.userKey)
	if err != nil {
		return "", "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	iv := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return "", "", err
	}

	encrypted := gcm.Seal(nil, iv, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(encrypted), base64.StdEncoding.EncodeToString(iv), nil
}

// DecryptData decrypts data on client side
func (ce *ClientEncryption) DecryptData(encryptedData string, iv string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(ce.userKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decrypted, err := gcm.Open(nil, ivBytes, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
