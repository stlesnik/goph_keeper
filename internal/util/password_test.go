package util

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "Password123",
			wantErr:  false,
		},
		{
			name:     "too short",
			password: "Pass1",
			wantErr:  true,
		},
		{
			name:     "too long",
			password: "Pass1234567890123456789012345678901234567890123456789012345678901234567890",
			wantErr:  true,
		},
		{
			name:     "no uppercase",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "no lowercase",
			password: "PASSWORD123",
			wantErr:  true,
		},
		{
			name:     "no digits",
			password: "Password",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "TestPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}

	if hash == password {
		t.Error("HashPassword() returned plain password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "TestPassword123"
	wrongPassword := "WrongPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}

	err = CheckPassword(password, hash)
	if err != nil {
		t.Errorf("CheckPassword() with correct password failed: %v", err)
	}

	err = CheckPassword(wrongPassword, hash)
	if err == nil {
		t.Error("CheckPassword() with wrong password should fail")
	}
}
