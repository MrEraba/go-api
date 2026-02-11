package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ivan-almanza/notes-api/internal/store"
)

// MockNoteStore implements store.NoteStorer for testing
type MockNoteStore struct {
	CreateNoteFunc func(ctx context.Context, note *store.Note) error
	ListNotesFunc  func(ctx context.Context, userID string) ([]*store.Note, error)
}

func (m *MockNoteStore) CreateNote(ctx context.Context, note *store.Note) error {
	if m.CreateNoteFunc != nil {
		return m.CreateNoteFunc(ctx, note)
	}
	return nil
}

func (m *MockNoteStore) ListNotes(ctx context.Context, userID string) ([]*store.Note, error) {
	if m.ListNotesFunc != nil {
		return m.ListNotesFunc(ctx, userID)
	}
	return nil, nil
}

func TestCreateNote_Authorized(t *testing.T) {
	userID := "user-123"
	mockStore := &MockNoteStore{
		CreateNoteFunc: func(ctx context.Context, note *store.Note) error {
			if note.UserID != userID {
				t.Errorf("Expected UserID %v, got %v", userID, note.UserID)
			}
			note.ID = "new-note-id"
			return nil
		},
	}
	handler := NewNotesHandler(mockStore)

	payload := map[string]string{
		"content": "My secret note",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(body))

	// Simulate middleware injecting UserID
	ctx := context.WithValue(req.Context(), ContextKeyUserID, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.CreateNote(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["id"] != "new-note-id" {
		t.Error("Expected ID in response")
	}
}

func TestGetNotes_Format(t *testing.T) {
	userID := "user-123"
	mockStore := &MockNoteStore{
		ListNotesFunc: func(ctx context.Context, uid string) ([]*store.Note, error) {
			if uid != userID {
				t.Errorf("Expected UserID %v, got %v", userID, uid)
			}
			return []*store.Note{
				{ID: "1", Content: "Note 1", CreatedAt: time.Now()},
				{ID: "2", Content: "Note 2", CreatedAt: time.Now()},
			}, nil
		},
	}
	handler := NewNotesHandler(mockStore)

	req := httptest.NewRequest(http.MethodGet, "/notes", nil)

	// Simulate middleware injecting UserID
	ctx := context.WithValue(req.Context(), ContextKeyUserID, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.GetNotes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := response["data"]; !ok {
		t.Error("Response JSON should contain 'data' field")
	}

	if _, ok := response["meta"]; !ok {
		t.Error("Response JSON should contain 'meta' field")
	}

	data := response["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(data))
	}
}
