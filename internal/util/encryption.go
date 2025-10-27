package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

// EncryptionService handles encryption and decryption of sensitive data
type EncryptionService struct {
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

// EncryptData encrypts data using AES-256-GCM
func (es *EncryptionService) EncryptData(data string, userKey []byte) (string, string, error) {
	block, err := aes.NewCipher(userKey)
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

// DecryptData decrypts data using AES-256-GCM
func (es *EncryptionService) DecryptData(encryptedData string, iv string, userKey []byte) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(userKey)
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

// DeriveUserKey derives a user key from password using PBKDF2
func (es *EncryptionService) DeriveUserKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}

// GenerateSalt generates a random salt
func (es *EncryptionService) GenerateSalt() ([]byte, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
