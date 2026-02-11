package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ivan-almanza/notes-api/internal/auth"
)

func TestAuthMiddleware_NoHeader(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Dummy handler should not be executed")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	WithAuth(dummyHandler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_BadFormat(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Dummy handler should not be executed")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token xyz")
	w := httptest.NewRecorder()

	WithAuth(dummyHandler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Dummy handler should not be executed")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	WithAuth(dummyHandler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_Success(t *testing.T) {
	userID := "user-123"
	token, _ := auth.GenerateToken(userID)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUserID := r.Context().Value(ContextKeyUserID)
		if ctxUserID == nil {
			t.Error("Context userID is nil")
		}
		if ctxUserID != userID {
			t.Errorf("Expected userID %v, got %v", userID, ctxUserID)
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	WithAuth(dummyHandler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
