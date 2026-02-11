package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_Structure(t *testing.T) {
	password := "mysecretpassword"
	hash, err := Hash(password)

	if err != nil {
		t.Fatalf("Hash() returned error: %v", err)
	}

	if hash == "" {
		t.Error("Hash() returned empty string")
	}

	if hash == password {
		t.Error("Hash() returned the original password")
	}
}

func TestComparePassword_Success(t *testing.T) {
	password := "password123"
	hash, err := Hash(password)
	if err != nil {
		t.Fatalf("Setup failed: Hash() returned error: %v", err)
	}

	err = Compare(password, hash)
	if err != nil {
		t.Errorf("Compare() failed for correct password: %v", err)
	}
}

func TestComparePassword_Failure(t *testing.T) {
	password := "password123"
	hash, err := Hash(password)
	if err != nil {
		t.Fatalf("Setup failed: Hash() returned error: %v", err)
	}

	err = Compare("wrongpassword", hash)
	if err != bcrypt.ErrMismatchedHashAndPassword {
		t.Errorf("Compare() expected ErrMismatchedHashAndPassword, got %v", err)
	}
}
