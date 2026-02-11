package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var testSecret = []byte("test-secret")

func init() {
	Secret = testSecret
}

func TestGenerateToken_ContainsClaims(t *testing.T) {
	userID := "user-123"
	tokenString, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return testSecret, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to parse claims")
	}

	if sub, ok := claims["sub"].(string); !ok || sub != userID {
		t.Errorf("Expected sub %v, got %v", userID, sub)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("Expiration claim missing")
	}

	if time.Unix(int64(exp), 0).Before(time.Now()) {
		t.Error("Token already expired")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	userID := "user-123"
	tokenString, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	token, err := ValidateToken(tokenString)
	if err != nil {
		t.Errorf("ValidateToken failed for valid token: %v", err)
	}

	if !token.Valid {
		t.Error("Token should be valid")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	// Manually create an expired token
	claims := jwt.MapClaims{
		"sub": "user-expired",
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(testSecret)

	_, err := ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestValidateToken_TamperedSignature(t *testing.T) {
	// Sign with a different secret
	claims := jwt.MapClaims{
		"sub": "user-tampered",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("wrong-secret"))

	_, err := ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for tampered token, got nil")
	}
}
