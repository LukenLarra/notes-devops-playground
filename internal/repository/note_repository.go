package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luken/notes-devops-playground/internal/model"
)

var ErrNoteNotFound = errors.New("note not found")

type NoteRepository interface {
	Create(ctx context.Context, title, content string) (model.Note, error)
	GetByID(ctx context.Context, id string) (model.Note, error)
	GetAll(ctx context.Context) ([]model.Note, error)
	Update(ctx context.Context, id, title, content string) (model.Note, error)
	Delete(ctx context.Context, id string) error
}

type noteRepository struct {
	pool *pgxpool.Pool
}

func NewNoteRepository(pool *pgxpool.Pool) NoteRepository {
	return &noteRepository{pool: pool}
}

func (r *noteRepository) Create(ctx context.Context, title, content string) (model.Note, error) {
	var note model.Note
	err := r.pool.QueryRow(ctx,
		`INSERT INTO notes (title, content) VALUES ($1, $2)
         RETURNING id, title, content, created_at, updated_at`,
		title, content,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	return note, err
}

func (r *noteRepository) GetByID(ctx context.Context, id string) (model.Note, error) {
	var note model.Note
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, content, created_at, updated_at FROM notes WHERE id = $1`,
		id,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err == pgx.ErrNoRows {
		return note, ErrNoteNotFound
	}
	return note, err
}

func (r *noteRepository) GetAll(ctx context.Context) ([]model.Note, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, title, content, created_at, updated_at FROM notes ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []model.Note
	for rows.Next() {
		var note model.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, rows.Err()
}

func (r *noteRepository) Update(ctx context.Context, id, title, content string) (model.Note, error) {
	var note model.Note
	err := r.pool.QueryRow(ctx,
		`UPDATE notes SET title = $1, content = $2, updated_at = NOW()
         WHERE id = $3
         RETURNING id, title, content, created_at, updated_at`,
		title, content, id,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err == pgx.ErrNoRows {
		return note, ErrNoteNotFound
	}
	return note, err
}

func (r *noteRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoteNotFound
	}
	return nil
}
