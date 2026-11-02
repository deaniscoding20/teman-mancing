package security

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenService(t *testing.T) {
	service := NewTokenService("test-secret-key-at-least-32-chars-long", 1*time.Hour)
	userID := uuid.New()
	username := "anglerinspire"

	token, err := service.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	if token == "" {
		t.Fatalf("generated token is empty")
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("expected username %v, got %v", username, claims.Username)
	}
}

func TestTokenServiceInvalidKey(t *testing.T) {
	service1 := NewTokenService("secret-one-at-least-32-chars-long", 1*time.Hour)
	service2 := NewTokenService("secret-two-at-least-32-chars-long", 1*time.Hour)

	token, _ := service1.GenerateToken(uuid.New(), "tester")
	_, err := service2.ValidateToken(token)
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken for mismatched secret, got %v", err)
	}
}
