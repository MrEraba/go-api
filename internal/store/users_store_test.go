package store

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func TestCreateUser_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %v", err)
	}
	defer db.Close()

	store := &PostgresStore{db: db}

	// UUID regex pattern for matching generated ID
	uuidPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at`)).
		WithArgs("test@example.com", "hashedpassword").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow("550e8400-e29b-41d4-a716-446655440000", time.Now()))

	user := &User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err = store.Create(context.Background(), user)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	if user.ID == "" {
		t.Error("User ID should not be empty after creation")
	}

	matched, _ := regexp.MatchString(uuidPattern, user.ID)
	if !matched && user.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected UUID, got %v", user.ID)
	}

	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %v", err)
	}
	defer db.Close()

	store := &PostgresStore{db: db}

	// Create a postgres error for duplicate key
	pqError := &pq.Error{Code: "23505"} // unique_violation

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users`)).
		WithArgs("duplicate@example.com", "hashedpassword").
		WillReturnError(pqError)

	user := &User{
		Email:    "duplicate@example.com",
		Password: "hashedpassword",
	}

	err = store.Create(context.Background(), user)
	if err != ErrDuplicateEmail {
		t.Errorf("Expected ErrDuplicateEmail, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByEmail_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %v", err)
	}
	defer db.Close()

	store := &PostgresStore{db: db}

	expectedID := "550e8400-e29b-41d4-a716-446655440000"
	expectedTime := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password, created_at FROM users WHERE email = $1`)).
		WithArgs("found@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "created_at"}).
			AddRow(expectedID, "found@example.com", "hashedpassword", expectedTime))

	user, err := store.GetByEmail(context.Background(), "found@example.com")
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}

	if user.ID != expectedID {
		t.Errorf("Expected ID %v, got %v", expectedID, user.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db: %v", err)
	}
	defer db.Close()

	store := &PostgresStore{db: db}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password, created_at FROM users WHERE email = $1`)).
		WithArgs("ghost@user.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "created_at"})) // Empty result

	_, err = store.GetByEmail(context.Background(), "ghost@user.com")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
