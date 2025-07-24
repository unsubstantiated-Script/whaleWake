package util

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain-text password using bcrypt.
// Returns the hashed password as a string, or an error if hashing fails.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a plain-text password with a bcrypt hashed password.
// Returns true if the password matches the hash, false otherwise.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
