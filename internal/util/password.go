package util

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// ValidatePassword validates the provided password
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("too short password")
	}
	if len(password) > 72 {
		return errors.New("too long password")
	}

	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"

	err := ""
	if !strings.ContainsAny(password, upper) {
		err = err + "doesn't have upper case character; "
	}
	if !strings.ContainsAny(password, lower) {
		err = err + "doesn't have lower case character; "
	}
	if !strings.ContainsAny(password, digits) {
		err = err + "doesn't contain digits; "
	}
	if err != "" {
		err = err[:(len(err) - 2)]
		return errors.New(err)
	}
	return nil
}

// HashPassword generates hash of the user's password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword checks if provided password matches the user's one
func CheckPassword(password string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
