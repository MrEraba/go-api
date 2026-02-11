package store

import "testing"

func TestUser_Validate_Success(t *testing.T) {
	u := &User{
		Email:    "test@example.com",
		Password: "password123",
	}

	if err := u.Validate(); err != nil {
		t.Errorf("Expected nil error for valid user, got %v", err)
	}
}

func TestUser_Validate_BadEmail(t *testing.T) {
	u := &User{
		Email:    "invalid-email",
		Password: "password123",
	}

	if err := u.Validate(); err == nil {
		t.Error("Expected error for invalid email, got nil")
	}
}

func TestUser_Validate_ShortPassword(t *testing.T) {
	u := &User{
		Email:    "test@example.com",
		Password: "123",
	}

	if err := u.Validate(); err == nil {
		t.Error("Expected error for short password, got nil")
	}
}
