package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/luken/notes-devops-playground/internal/model"
	"github.com/luken/notes-devops-playground/internal/repository"
	"github.com/luken/notes-devops-playground/internal/service"
)

type mockRepo struct {
	notes []model.Note
	err   error
}

func (m *mockRepo) Create(ctx context.Context, title, content string) (model.Note, error) {
	if m.err != nil {
		return model.Note{}, m.err
	}
	n := model.Note{ID: "1", Title: title, Content: content}
	m.notes = append(m.notes, n)
	return n, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (model.Note, error) {
	if m.err != nil {
		return model.Note{}, m.err
	}
	for _, n := range m.notes {
		if n.ID == id {
			return n, nil
		}
	}
	return model.Note{}, repository.ErrNoteNotFound
}

func (m *mockRepo) GetAll(ctx context.Context) ([]model.Note, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.notes, nil
}

func (m *mockRepo) Update(ctx context.Context, id, title, content string) (model.Note, error) {
	if m.err != nil {
		return model.Note{}, m.err
	}
	for i, n := range m.notes {
		if n.ID == id {
			m.notes[i].Title = title
			m.notes[i].Content = content
			return m.notes[i], nil
		}
	}
	return model.Note{}, repository.ErrNoteNotFound
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	for i, n := range m.notes {
		if n.ID == id {
			m.notes = append(m.notes[:i], m.notes[i+1:]...)
			return nil
		}
	}
	return repository.ErrNoteNotFound
}

type mockValidator struct{}

func (m *mockValidator) ValidateCreate(req model.CreateNoteRequest) error {
	return nil
}

func (m *mockValidator) ValidateUpdate(req model.UpdateNoteRequest) error {
	return nil
}

func TestNoteService_Create(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewNoteService(repo, &mockValidator{})

	note, err := svc.Create(context.Background(), "Title", "Content")
	if err != nil {
		t.Fatal(err)
	}
	if note.Title != "Title" {
		t.Errorf("expected 'Title', got %q", note.Title)
	}
	if len(repo.notes) != 1 {
		t.Errorf("expected repo to have 1 note, got %d", len(repo.notes))
	}
}

func TestNoteService_Create_EmptyTitle(t *testing.T) {
	svc := service.NewNoteService(&mockRepo{}, &service.DefaultValidator{})

	_, err := svc.Create(context.Background(), "", "Content")
	if err != service.ErrValidation {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestNoteService_GetByID_NotFound(t *testing.T) {
	svc := service.NewNoteService(&mockRepo{}, &mockValidator{})

	_, err := svc.GetByID(context.Background(), "nonexistent")
	if err != service.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestNoteService_Delete(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewNoteService(repo, &mockValidator{})

	note, err := svc.Create(context.Background(), "Delete Me", "Content")
	if err != nil {
		t.Fatal(err)
	}

	err = svc.Delete(context.Background(), note.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.GetByID(context.Background(), note.ID)
	if err != service.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestNoteService_Update(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewNoteService(repo, &mockValidator{})

	note, err := svc.Create(context.Background(), "Original", "Original content")
	if err != nil {
		t.Fatal(err)
	}

	updated, err := svc.Update(context.Background(), note.ID, "Updated", "Updated content")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "Updated" {
		t.Errorf("expected 'Updated', got %q", updated.Title)
	}
}

func TestNoteService_GetAll(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewNoteService(repo, &mockValidator{})

	svc.Create(context.Background(), "Note 1", "Content 1")
	svc.Create(context.Background(), "Note 2", "Content 2")

	notes, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 2 {
		t.Errorf("expected 2 notes, got %d", len(notes))
	}
}

func TestNoteService_RepoError(t *testing.T) {
	repo := &mockRepo{err: errors.New("db connection failed")}
	svc := service.NewNoteService(repo, &mockValidator{})

	_, err := svc.Create(context.Background(), "Title", "Content")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "db connection failed" {
		t.Errorf("expected 'db connection failed', got %v", err)
	}
}
