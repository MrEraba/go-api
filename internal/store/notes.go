package store

import (
	"context"
	"time"
)

type Note struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteStorer interface {
	CreateNote(ctx context.Context, note *Note) error
	ListNotes(ctx context.Context, userID string) ([]*Note, error)
}

func (s *PostgresStore) CreateNote(ctx context.Context, note *Note) error {
	query := `INSERT INTO notes (user_id, content) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query, note.UserID, note.Content).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) ListNotes(ctx context.Context, userID string) ([]*Note, error) {
	query := `SELECT id, user_id, content, created_at, updated_at FROM notes WHERE user_id = $1`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*Note
	for rows.Next() {
		note := &Note{}
		if err := rows.Scan(&note.ID, &note.UserID, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
