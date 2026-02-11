package store

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/lib/pq"
)

var (
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
	ErrDuplicateEmail   = errors.New("email already exists")
	ErrNotFound         = errors.New("resource not found")
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"` // plaintext for input, not stored
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	if len(u.Password) < 6 {
		return ErrPasswordTooShort
	}

	// Simple regex for email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if match, _ := regexp.MatchString(emailRegex, u.Email); !match {
		return ErrInvalidEmail
	}

	return nil
}

type UserStorer interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (s *PostgresStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`

	var user User
	err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
