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

// Response is a struct for HTTP response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// EncryptedDataItem represents an encrypted data item stored by user
type EncryptedDataItem struct {
	ID            string
	UserID        string
	Type          string
	Title         string
	EncryptedData string
	IV            string
	Metadata      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
