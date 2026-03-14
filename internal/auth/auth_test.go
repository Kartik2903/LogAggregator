package auth

// Unit tests for the bcrypt authentication module.
// LAB 8: Table-driven tests, testing security-critical code.

import (
	"testing"
)

// Test user registration and authentication with bcrypt.
func TestRegisterAndAuthenticate(t *testing.T) {
	store := NewUserStore()

	// Register a user — password gets bcrypt-hashed internally
	err := store.RegisterUser("admin", "securePassword123")
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	// LAB 8: Table-driven tests for authentication
	tests := []struct {
		name     string
		username string
		password string
		want     bool
	}{
		{"correct credentials", "admin", "securePassword123", true},
		{"wrong password", "admin", "wrongPassword", false},
		{"non-existent user", "unknown", "securePassword123", false},
		{"empty password", "admin", "", false},
		{"case-sensitive password", "admin", "SecurePassword123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := store.Authenticate(tt.username, tt.password)
			if got != tt.want {
				t.Errorf("Authenticate(%q, %q) = %v, want %v",
					tt.username, tt.password, got, tt.want)
			}
		})
	}
}

// Test that duplicate registration is rejected.
func TestDuplicateRegistration(t *testing.T) {
	store := NewUserStore()

	err := store.RegisterUser("admin", "password1")
	if err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	err = store.RegisterUser("admin", "password2")
	if err == nil {
		t.Error("expected error for duplicate registration, got nil")
	}
}

// Test that bcrypt produces different hashes for the same password.
// (This proves that bcrypt uses random salts.)
func TestBcryptUniqueSalts(t *testing.T) {
	store := NewUserStore()

	_ = store.RegisterUser("user1", "samePassword")
	_ = store.RegisterUser("user2", "samePassword")

	store.mu.RLock()
	hash1 := store.users["user1"]
	hash2 := store.users["user2"]
	store.mu.RUnlock()

	if hash1 == hash2 {
		t.Error("bcrypt should produce different hashes for same password (different salts)")
	}

	// Both should still authenticate correctly
	if !store.Authenticate("user1", "samePassword") {
		t.Error("user1 should authenticate")
	}
	if !store.Authenticate("user2", "samePassword") {
		t.Error("user2 should authenticate")
	}
}

// Test session management
func TestSessionManagement(t *testing.T) {
	sm := NewSessionManager()

	// Create session
	token, err := sm.CreateSession("admin")
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	if token == "" {
		t.Fatal("session token should not be empty")
	}

	// Validate session
	username, valid := sm.ValidateSession(token)
	if !valid {
		t.Error("session should be valid")
	}
	if username != "admin" {
		t.Errorf("expected username 'admin', got %q", username)
	}

	// Invalid token
	_, valid = sm.ValidateSession("bogus-token")
	if valid {
		t.Error("bogus token should not be valid")
	}

	// Destroy session (logout)
	sm.DestroySession(token)
	_, valid = sm.ValidateSession(token)
	if valid {
		t.Error("session should be invalid after destruction")
	}
}
