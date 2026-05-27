package handlers

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

const (
	CodeNotFound        = "NOT_FOUND"
	CodeInvalidJSON     = "INVALID_JSON"
	CodeValidationError = "VALIDATION_ERROR"
	CodeInternalError   = "INTERNAL_ERROR"
)

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{Error: message, Code: code})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
