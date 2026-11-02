package user

import (
	"testing"
)

func TestNewUserValid(t *testing.T) {
	u, err := NewUser("Budi Santoso", "budimancing", "budi@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error creating user: %v", err)
	}
	if u.Name != "Budi Santoso" {
		t.Errorf("expected name 'Budi Santoso', got '%s'", u.Name)
	}
	if u.Username != "budimancing" {
		t.Errorf("expected username 'budimancing', got '%s'", u.Username)
	}
	if !u.CheckPassword("password123") {
		t.Errorf("password verification failed for valid password")
	}
	if u.CheckPassword("wrongpassword") {
		t.Errorf("password verification succeeded for invalid password")
	}
}

func TestNewUserValidationErrors(t *testing.T) {
	cases := []struct {
		name     string
		username string
		email    string
		password string
		expected string
	}{
		{"", "testuser", "test@example.com", "password123", "name cannot be empty"},
		{"Valid Name", "ab", "test@example.com", "password123", ErrInvalidUsername.Error()},
		{"Valid Name", "testuser", "notanemail", "password123", ErrInvalidEmail.Error()},
		{"Valid Name", "testuser", "test@example.com", "short", ErrPasswordTooWeak.Error()},
	}

	for _, tc := range cases {
		_, err := NewUser(tc.name, tc.username, tc.email, tc.password)
		if err == nil {
			t.Errorf("expected error for case %v, got nil", tc)
		}
	}
}
