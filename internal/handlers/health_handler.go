package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	pool *pgxpool.Pool
}

func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	dbOk := true
	if err := h.pool.Ping(context.Background()); err != nil {
		dbOk = false
	}

	status := http.StatusOK
	statusText := "ok"
	if !dbOk {
		status = http.StatusServiceUnavailable
		statusText = "degraded"
	}

	writeJSON(w, status, map[string]interface{}{
		"status": statusText,
		"db":     dbOk,
	})
}
