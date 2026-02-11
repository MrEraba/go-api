package auth

import "golang.org/x/crypto/bcrypt"

// Hash generates a bcrypt hash of the password.
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare compares a bcrypt hashed password with its possible plaintext equivalent.
func Compare(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
