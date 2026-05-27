package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luken/notes-devops-playground/internal/handlers"
	"github.com/luken/notes-devops-playground/internal/repository"
	"github.com/luken/notes-devops-playground/internal/service"
)

func setupTestRouter(t *testing.T) *chi.Mux {
	t.Helper()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	repo := repository.NewNoteRepository(pool)
	svc := service.NewNoteService(repo, &service.DefaultValidator{})
	h := handlers.NewNoteHandler(svc)

	r := chi.NewRouter()
	r.Get("/health", handlers.HealthHandler())
	r.Get("/api/notes", h.GetAll)
	r.Post("/api/notes", h.Create)
	r.Get("/api/notes/{id}", h.GetByID)
	r.Put("/api/notes/{id}", h.Update)
	r.Delete("/api/notes/{id}", h.Delete)

	return r
}

func TestHealthEndpoint(t *testing.T) {
    router := chi.NewRouter()
    router.Get("/health", handlers.HealthHandler())
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateNote(t *testing.T) {
	router := setupTestRouter(t)

	body := map[string]string{"title": "Integration Test", "content": "Testing"}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/notes", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var note map[string]any
	json.NewDecoder(rec.Body).Decode(&note)
	if note["title"] != "Integration Test" {
		t.Errorf("expected 'Integration Test', got %v", note["title"])
	}
}

func TestCreateNote_EmptyTitle(t *testing.T) {
	router := setupTestRouter(t)

	body := map[string]string{"title": "", "content": "Some content"}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/notes", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestGetNonExistentNote(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/notes/00000000-0000-0000-0000-000000000000", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
