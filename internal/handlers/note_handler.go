package handlers

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/luken/notes-devops-playground/internal/service"
)

const maxBodySize = 1 << 20 // 1 MB

type noteHandler struct {
	svc service.NoteService
}

func NewNoteHandler(svc service.NoteService) *noteHandler {
	return &noteHandler{svc: svc}
}

func (h *noteHandler) Create(w http.ResponseWriter, r *http.Request) {
    r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

    var req struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeInvalidJSON, "Invalid request body")
		return
	}

	note, err := h.svc.Create(r.Context(), req.Title, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			writeError(w, http.StatusBadRequest, CodeValidationError, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, CodeInternalError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

func (h *noteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	note, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, CodeNotFound, "Note not found")
			return
		}
		writeError(w, http.StatusInternalServerError, CodeInternalError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *noteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes, err := h.svc.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, CodeInternalError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"notes": notes,
		"total": len(notes),
	})
}

func (h *noteHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

    var req struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeInvalidJSON, "Invalid request body")
		return
	}

	note, err := h.svc.Update(r.Context(), id, req.Title, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, CodeNotFound, "Note not found")
			return
		}
		if errors.Is(err, service.ErrValidation) {
			writeError(w, http.StatusBadRequest, CodeValidationError, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, CodeInternalError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *noteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, CodeNotFound, "Note not found")
			return
		}
		writeError(w, http.StatusInternalServerError, CodeInternalError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
