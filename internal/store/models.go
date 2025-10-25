package store

import "time"

// User is a struct that repeats db entity
type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
