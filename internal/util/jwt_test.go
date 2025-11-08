package util

import (
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

func TestGenerateAndParseJWT(t *testing.T) {
	userID := "test-user-id"
	username := "testuser"
	email := "test@example.com"
	secret := "test-secret-that-is-long-enough-for-validation"

	token, err := GenerateJWT(userID, username, email, secret)
	if err != nil {
		t.Fatalf("GenerateJWT() failed: %v", err)
	}

	if token == "" {
		t.Error("GenerateJWT() returned empty token")
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken() failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("ParseToken() userID = %v, want %v", claims.UserID, userID)
	}

	if claims.Username != username {
		t.Errorf("ParseToken() username = %v, want %v", claims.Username, username)
	}

	if claims.Email != email {
		t.Errorf("ParseToken() email = %v, want %v", claims.Email, email)
	}
}

func TestParseToken_Expired(t *testing.T) {
	secret := "test-secret-that-is-long-enough-for-validation"

	// Create expired token manually
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
		UserID:   "test-user",
		Username: "testuser",
		Email:    "test@example.com",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	_, err = ParseToken(tokenString, secret)
	if err == nil {
		t.Error("ParseToken() with expired token should fail")
	}
}

func TestParseToken_InvalidSignature(t *testing.T) {
	userID := "test-user-id"
	username := "testuser"
	email := "test@example.com"
	secret := "test-secret-that-is-long-enough-for-validation"
	wrongSecret := "wrong-secret-that-is-long-enough-for-validation"

	token, err := GenerateJWT(userID, username, email, secret)
	if err != nil {
		t.Fatalf("GenerateJWT() failed: %v", err)
	}

	_, err = ParseToken(token, wrongSecret)
	if err == nil {
		t.Error("ParseToken() with wrong secret should fail")
	}
}
