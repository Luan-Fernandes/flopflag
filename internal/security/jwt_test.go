package security

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenManager_SignAndParse(t *testing.T) {
	m := NewTokenManager("test-secret", 5*time.Minute)

	userID := uuid.New()
	token, err := m.Sign(userID, "test@example.com", "admin")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	gotID, gotRole, err := m.Parse(token)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if gotID != userID {
		t.Fatalf("userID mismatch: got %v want %v", gotID, userID)
	}
	if gotRole != "admin" {
		t.Fatalf("role mismatch: got %q want %q", gotRole, "admin")
	}
}
