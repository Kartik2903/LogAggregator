package auth

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Run with: go test -v -run TestPrintHash ./internal/auth/
func TestPrintHash(t *testing.T) {
	password := "mySecret123"

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	fmt.Println("\n  Password:   ", password)
	fmt.Println("  Bcrypt Hash:", string(hash))

	hash2, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	fmt.Println("  Hash Again: ", string(hash2))
	fmt.Println("  Same hash?  ", string(hash) == string(hash2))
}
