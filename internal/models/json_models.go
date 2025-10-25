package models

import "time"

// RegisterUserRequest is a struct for user registration request
type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginUserRequest is a struct for user log in request
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// DataItem represents a data item stored by user
type DataItem struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type" validate:"required,oneof=password text binary card"`
	Title     string `json:"title" validate:"required"`
	Data      string `json:"data" validate:"required"`
	Metadata  string `json:"metadata,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
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

// CreateDataRequest is a struct for creating new data item
type CreateDataRequest struct {
	Type     string `json:"type" validate:"required,oneof=password text binary card"`
	Title    string `json:"title" validate:"required"`
	Data     string `json:"data" validate:"required"`
	Metadata string `json:"metadata,omitempty"`
}

// UpdateDataRequest is a struct for updating data item
type UpdateDataRequest struct {
	Type     string `json:"type" validate:"required,oneof=password text binary card"`
	Title    string `json:"title" validate:"required"`
	Data     string `json:"data" validate:"required"`
	Metadata string `json:"metadata,omitempty"`
}

// ChangePasswordRequest is a struct for changing user password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// UserProfile represents user profile information
type UserProfile struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// VersionResponse represents server version information
type VersionResponse struct {
	Version string `json:"version"`
	Build   string `json:"build"`
	Date    string `json:"date"`
}
