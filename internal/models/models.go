package models

import "time"

// UserContext represents user data from jwt token
type UserContext struct {
	UserID string
	Email  string
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
