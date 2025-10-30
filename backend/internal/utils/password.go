package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt
// Bcrypt is slow ON PURPOSE - makes brute force attacks impractical
// Cost of 12 means ~250ms to hash (perfect balance)
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Bcrypt cost: 10 = fast (testing), 12 = recommended, 14 = paranoid
	const cost = 12

	// bcrypt automatically adds salt, so each hash is unique
	// Even same password produces different hashes!
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword checks if a password matches its hash
// Constant-time comparison prevents timing attacks
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
