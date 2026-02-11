package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ivan-almanza/notes-api/internal/auth"
	"github.com/ivan-almanza/notes-api/internal/store"
)

// MockUserStore implements store.UserStorer for testing
type MockUserStore struct {
	CreateFunc     func(ctx context.Context, user *store.User) error
	GetByEmailFunc func(ctx context.Context, email string) (*store.User, error)
}

func (m *MockUserStore) Create(ctx context.Context, user *store.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*store.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func TestRegister_Success(t *testing.T) {
	mockStore := &MockUserStore{
		CreateFunc: func(ctx context.Context, user *store.User) error {
			user.ID = "test-uuid"
			return nil
		},
	}
	handler := NewAuthHandler(mockStore)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "test-uuid") {
		t.Error("Response body should contain user ID")
	}
}

func TestRegister_InvalidInput(t *testing.T) {
	mockStore := &MockUserStore{
		CreateFunc: func(ctx context.Context, user *store.User) error {
			t.Error("Create should not be called")
			return nil
		},
	}
	handler := NewAuthHandler(mockStore)

	payload := map[string]string{
		"email":    "invalid-email",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, got %d", w.Code)
	}
}

func TestRegister_Duplicate(t *testing.T) {
	mockStore := &MockUserStore{
		CreateFunc: func(ctx context.Context, user *store.User) error {
			return store.ErrDuplicateEmail
		},
	}
	handler := NewAuthHandler(mockStore)

	payload := map[string]string{
		"email":    "duplicate@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 Conflict, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	hashedPassword, _ := auth.Hash("password123")
	mockStore := &MockUserStore{
		GetByEmailFunc: func(ctx context.Context, email string) (*store.User, error) {
			return &store.User{
				ID:       "user-123",
				Email:    "test@example.com",
				Password: hashedPassword,
			}, nil
		},
	}
	handler := NewAuthHandler(mockStore)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response")
	}

	if _, ok := response["token"]; !ok {
		t.Error("Response should contain token")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	mockStore := &MockUserStore{
		GetByEmailFunc: func(ctx context.Context, email string) (*store.User, error) {
			return nil, store.ErrNotFound
		},
	}
	handler := NewAuthHandler(mockStore)

	payload := map[string]string{
		"email":    "ghost@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", w.Code)
	}
}

func TestLogin_BadPassword(t *testing.T) {
	hashedPassword, _ := auth.Hash("password123")
	mockStore := &MockUserStore{
		GetByEmailFunc: func(ctx context.Context, email string) (*store.User, error) {
			return &store.User{
				ID:       "user-123",
				Email:    "test@example.com",
				Password: hashedPassword,
			}, nil
		},
	}
	handler := NewAuthHandler(mockStore)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", w.Code)
	}
}
