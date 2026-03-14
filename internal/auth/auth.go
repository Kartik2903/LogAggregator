package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// Default bcrypt cost factor.
// Higher cost = slower hashing = harder to brute-force.
// Cost 12 takes ~250ms per hash on modern hardware — good balance.
const DefaultCost = 12

// ============================================================
// UserStore manages user credentials with bcrypt-hashed passwords.
// In a real application, this would be backed by a database.
// Here we use an in-memory map for demonstration.
// ============================================================
type UserStore struct {
	// Map of username → bcrypt-hashed password
	users map[string]string // LAB 4: Map usage
	mu    sync.RWMutex      // LAB 9: Mutex for concurrent safety
}

// NewUserStore creates a new empty user store.
func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]string),
	}
}

// ============================================================
// RegisterUser hashes the plaintext password with bcrypt and stores it.
//
// BCRYPT FLOW:
//
//	plaintext "mypassword"
//	  → bcrypt.GenerateFromPassword()
//	    → adds random salt + hashes with cost factor
//	      → "$2a$12$LJ3m4ys..." (60-char hash string)
//
// The hash includes: algorithm version, cost, salt, and hash — all in one string.
// ============================================================
func (us *UserStore) RegisterUser(username, plainPassword string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	// Check if user already exists
	if _, exists := us.users[username]; exists {
		return fmt.Errorf("user %q already exists", username)
	}

	// BCRYPT: Hash the plaintext password.
	// GenerateFromPassword(password, cost) → hashed bytes
	// The salt is generated automatically — no need to manage it ourselves.
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword), // plaintext password as bytes
		DefaultCost,           // cost factor (2^12 iterations)
	)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Store the hash — never store the plaintext password!
	us.users[username] = string(hashedBytes)
	return nil
}

// ============================================================
// Authenticate verifies a plaintext password against the stored bcrypt hash.
//
// BCRYPT COMPARISON:
//
//	bcrypt.CompareHashAndPassword(storedHash, attemptedPassword)
//	  → extracts the salt from the stored hash
//	  → hashes the attempted password with the same salt and cost
//	  → compares the two hashes in constant time (prevents timing attacks)
//	  → returns nil if they match, error if they don't
//
// ============================================================
func (us *UserStore) Authenticate(username, plainPassword string) bool {
	us.mu.RLock()
	defer us.mu.RUnlock()

	storedHash, exists := us.users[username]
	if !exists {
		// User not found — still run bcrypt to prevent timing attacks
		// (an attacker can't tell if the username was wrong vs password)
		bcrypt.CompareHashAndPassword(
			[]byte("$2a$12$dummyhashtopreventtimingattacks000000000000000000"),
			[]byte(plainPassword),
		)
		return false
	}

	// BCRYPT: Compare the stored hash with the attempted password.
	// Returns nil on success, error on mismatch.
	err := bcrypt.CompareHashAndPassword(
		[]byte(storedHash),    // the stored bcrypt hash
		[]byte(plainPassword), // the attempted plaintext password
	)

	return err == nil // true = password matches, false = wrong password
}

// ============================================================
// Session Management — simple token-based sessions.
// After successful login, a random session token is generated.
// The token is stored server-side and sent to the client as a cookie.
// ============================================================

// SessionManager handles active session tokens.
type SessionManager struct {
	sessions map[string]string // token → username
	mu       sync.RWMutex
}

// NewSessionManager creates a new session manager.
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]string),
	}
}

// CreateSession generates a random token and associates it with a username.
func (sm *SessionManager) CreateSession(username string) (string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Generate a cryptographically secure random token (32 bytes = 64 hex chars)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	token := hex.EncodeToString(tokenBytes)
	sm.sessions[token] = username
	return token, nil
}

// ValidateSession checks if a token is valid and returns the username.
func (sm *SessionManager) ValidateSession(token string) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	username, exists := sm.sessions[token]
	return username, exists
}

// DestroySession removes a session (logout).
func (sm *SessionManager) DestroySession(token string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, token)
}
