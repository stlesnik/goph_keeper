package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenExp is JWT expiration.
const TokenExp = time.Hour * 24

// Claims is JWT claims with user ID.
type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GenerateJWT creates a signed JWT for the given user.
func GenerateJWT(username, email string, secretKey string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		Username: username,
		Email:    email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "nil", err
	}

	return tokenString, nil
}

// ParseToken takes user's token from the HTTP headers, validates it and returns claims.
func ParseToken(signedToken string, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}
	return claims, nil
}
