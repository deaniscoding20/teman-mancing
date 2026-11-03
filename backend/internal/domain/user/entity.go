package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrInvalidUsername = errors.New("username must be 3-30 characters alphanumeric")
	ErrPasswordTooWeak = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong = errors.New("password cannot exceed 72 characters")
)

// User represents a user account and profile in Teman Mancing.
type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AvatarKey    *string   `json:"avatar_key,omitempty"`
	Bio          *string   `json:"bio,omitempty"`
	Location     *string   `json:"location,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser creates a new User with validated fields and hashed password.
func NewUser(name, username, email, rawPassword string) (*User, error) {
	name = strings.TrimSpace(name)
	username = strings.ToLower(strings.TrimSpace(username))
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if len(username) < 3 || len(username) > 30 {
		return nil, ErrInvalidUsername
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, ErrInvalidEmail
	}
	if len(rawPassword) < 8 {
		return nil, ErrPasswordTooWeak
	}
	if len(rawPassword) > 72 {
		return nil, ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		ID:           uuid.New(),
		Name:         name,
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// CheckPassword verifies if raw password matches the stored bcrypt hash.
func (u *User) CheckPassword(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(rawPassword))
	return err == nil
}

// SetPassword updates the password hash.
func (u *User) SetPassword(newPassword string) error {
	if len(newPassword) < 8 {
		return ErrPasswordTooWeak
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	u.UpdatedAt = time.Now().UTC()
	return nil
}
