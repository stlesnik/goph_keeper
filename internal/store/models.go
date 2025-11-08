package store

import "time"

// User is a struct that repeats db entity
type User struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Salt         string    `db:"salt"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// EncryptedDataItem represents an encrypted data item stored by user
type EncryptedDataItem struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	Type          string    `db:"type"`
	Title         string    `db:"title"`
	EncryptedData string    `db:"encrypted_data"`
	IV            string    `db:"iv"`
	Metadata      string    `db:"metadata"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
