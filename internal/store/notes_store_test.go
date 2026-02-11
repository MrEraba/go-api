package store

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateNote(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %v", err)
	}
	defer db.Close()

	store := &PostgresStore{db: db}

	userID := "user-123"
	noteContent := "This is a test note"

	// Mock INSERT
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO notes (user_id, content) VALUES ($1, $2) RETURNING id, created_at, updated_at`)).
		WithArgs(userID, noteContent).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("note-uuid", time.Now(), time.Now()))

	note := &Note{
		UserID:  userID,
		Content: noteContent,
	}

	err = store.CreateNote(context.Background(), note)
	if err != nil {
		t.Errorf("CreateNote failed: %v", err)
	}

	if note.ID == "" {
		t.Error("Note ID should not be empty")
	}
	if note.UserID != userID {
		t.Errorf("Expected UserID %v, got %v", userID, note.UserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListNotes_Isolation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %v", err)
	}
	defer db.Close()

	store := &PostgresStore{db: db}

	userID := "user-A"
	expectedContent := "User A Note"

	// Mock SELECT with WHERE user_id = $1
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, content, created_at, updated_at FROM notes WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content", "created_at", "updated_at"}).
			AddRow("note-1", userID, expectedContent, time.Now(), time.Now()))

	notes, err := store.ListNotes(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListNotes failed: %v", err)
	}

	if len(notes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(notes))
	}

	if notes[0].UserID != userID {
		t.Errorf("Expected note belong to %v, got %v", userID, notes[0].UserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
