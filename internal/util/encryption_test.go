package util

import (
	"testing"
)

func TestEncryptionService_EncryptDecryptData(t *testing.T) {
	es := NewEncryptionService()

	originalData := "sensitive data that needs encryption"
	password := "test-password"

	salt, err := es.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() failed: %v", err)
	}

	userKey := es.DeriveUserKey(password, salt)

	encryptedData, iv, err := es.EncryptData(originalData, userKey)
	if err != nil {
		t.Fatalf("EncryptData() failed: %v", err)
	}

	if encryptedData == originalData {
		t.Error("EncryptData() returned original data")
	}

	decryptedData, err := es.DecryptData(encryptedData, iv, userKey)
	if err != nil {
		t.Fatalf("DecryptData() failed: %v", err)
	}

	if decryptedData != originalData {
		t.Errorf("DecryptData() = %v, want %v", decryptedData, originalData)
	}
}

func TestEncryptionService_DeriveUserKey(t *testing.T) {
	es := NewEncryptionService()

	password := "test-password"
	salt := []byte("test-salt-that-is-32-bytes-long!")

	key1 := es.DeriveUserKey(password, salt)
	key2 := es.DeriveUserKey(password, salt)

	if len(key1) != 32 {
		t.Errorf("DeriveUserKey() key length = %v, want 32", len(key1))
	}

	// Same password and salt should produce same key
	for i := range key1 {
		if key1[i] != key2[i] {
			t.Error("DeriveUserKey() should be deterministic")
			break
		}
	}

	// Different salt should produce different key
	differentSalt := []byte("different-salt-that-is-32-bytes!")
	key3 := es.DeriveUserKey(password, differentSalt)

	different := false
	for i := range key1 {
		if key1[i] != key3[i] {
			different = true
			break
		}
	}

	if !different {
		t.Error("DeriveUserKey() with different salt should produce different key")
	}
}
