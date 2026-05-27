package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luken/notes-devops-playground/internal/repository"
)

func TestNoteRepository_CRUD(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewNoteRepository(pool)

	t.Run("create and get by id", func(t *testing.T) {
		note, err := repo.Create(ctx, "Test Title", "Test Content")
		if err != nil {
			t.Fatal(err)
		}
		if note.Title != "Test Title" {
			t.Errorf("expected 'Test Title', got %q", note.Title)
		}
		if note.Content != "Test Content" {
			t.Errorf("expected 'Test Content', got %q", note.Content)
		}
		if note.ID == "" {
			t.Error("expected non-empty ID")
		}

		got, err := repo.GetByID(ctx, note.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.Title != note.Title {
			t.Errorf("expected %q, got %q", note.Title, got.Title)
		}
	})

	t.Run("get all returns created notes", func(t *testing.T) {
		notes, err := repo.GetAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(notes) == 0 {
			t.Error("expected at least one note")
		}
	})

	t.Run("update note", func(t *testing.T) {
		note, err := repo.Create(ctx, "Original", "Original content")
		if err != nil {
			t.Fatal(err)
		}

		updated, err := repo.Update(ctx, note.ID, "Updated", "Updated content")
		if err != nil {
			t.Fatal(err)
		}
		if updated.Title != "Updated" {
			t.Errorf("expected 'Updated', got %q", updated.Title)
		}
	})

	t.Run("delete note", func(t *testing.T) {
		note, err := repo.Create(ctx, "Delete Me", "To be deleted")
		if err != nil {
			t.Fatal(err)
		}

		err = repo.Delete(ctx, note.ID)
		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.GetByID(ctx, note.ID)
		if err != repository.ErrNoteNotFound {
			t.Errorf("expected ErrNoteNotFound, got %v", err)
		}
	})

	t.Run("get non-existent returns error", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
		if err != repository.ErrNoteNotFound {
			t.Errorf("expected ErrNoteNotFound, got %v", err)
		}
	})
}
