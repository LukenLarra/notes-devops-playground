package service

import (
	"context"
	"errors"
	"strings"

	"github.com/luken/notes-devops-playground/internal/model"
	"github.com/luken/notes-devops-playground/internal/repository"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation error")
)

type NoteValidator interface {
	ValidateCreate(req model.CreateNoteRequest) error
	ValidateUpdate(req model.UpdateNoteRequest) error
}

type DefaultValidator struct{}

func (v *DefaultValidator) ValidateCreate(req model.CreateNoteRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return ErrValidation
	}
	if strings.TrimSpace(req.Content) == "" {
		return ErrValidation
	}
	return nil
}

func (v *DefaultValidator) ValidateUpdate(req model.UpdateNoteRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return ErrValidation
	}
	if strings.TrimSpace(req.Content) == "" {
		return ErrValidation
	}
	return nil
}

type NoteService interface {
	Create(ctx context.Context, title, content string) (model.Note, error)
	GetByID(ctx context.Context, id string) (model.Note, error)
	GetAll(ctx context.Context) ([]model.Note, error)
	Update(ctx context.Context, id, title, content string) (model.Note, error)
	Delete(ctx context.Context, id string) error
}

type noteService struct {
	repo      repository.NoteRepository
	validator NoteValidator
}

func NewNoteService(repo repository.NoteRepository, validator NoteValidator) NoteService {
	return &noteService{repo: repo, validator: validator}
}

func (s *noteService) Create(ctx context.Context, title, content string) (model.Note, error) {
	if err := s.validator.ValidateCreate(model.CreateNoteRequest{Title: title, Content: content}); err != nil {
		return model.Note{}, err
	}
	return s.repo.Create(ctx, title, content)
}

func (s *noteService) GetByID(ctx context.Context, id string) (model.Note, error) {
	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoteNotFound) {
			return model.Note{}, ErrNotFound
		}
		return model.Note{}, err
	}
	return note, nil
}

func (s *noteService) GetAll(ctx context.Context) ([]model.Note, error) {
	return s.repo.GetAll(ctx)
}

func (s *noteService) Update(ctx context.Context, id, title, content string) (model.Note, error) {
	if err := s.validator.ValidateUpdate(model.UpdateNoteRequest{Title: title, Content: content}); err != nil {
		return model.Note{}, err
	}
	note, err := s.repo.Update(ctx, id, title, content)
	if err != nil {
		if errors.Is(err, repository.ErrNoteNotFound) {
			return model.Note{}, ErrNotFound
		}
		return model.Note{}, err
	}
	return note, nil
}

func (s *noteService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoteNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
