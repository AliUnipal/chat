package httpsrv

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type ErrorType string

const (
	UnknownError        ErrorType = ""
	ValidationError     ErrorType = "validation_error"
	AuthenticationError ErrorType = "authentication_error"
	AuthorizationError  ErrorType = "authorization_error"
)

func httpCode(ctx context.Context, typ ErrorType) int {
	switch typ {
	case UnknownError:
		return http.StatusInternalServerError
	case ValidationError:
		return http.StatusBadRequest
	default:
		slog.ErrorContext(ctx, "unknown error type", "type", typ)
		return http.StatusInternalServerError
	}
}

type errorResponse[T any] struct {
	Type    ErrorType          `json:"type,omitempty"`
	Code    string             `json:"code,omitempty"`
	Message string             `json:"message,omitempty"`
	Data    T                  `json:"data,omitzero"`
	Errors  []errorResponse[T] `json:"errors,omitempty"`
}

func RespondWithError(ctx context.Context, w http.ResponseWriter, code string, err error) {
	respond(ctx, w, httpCode(ctx, UnknownError), errorResponse[struct{}]{
		Type:    UnknownError,
		Code:    code,
		Message: err.Error(),
	})
}

type validationError struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

type kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func KV(key, value string) kv {
	return kv{Key: key, Value: value}
}

func RespondWithValidationError(ctx context.Context, w http.ResponseWriter, code string, values ...kv) {
	var errors []validationError
	for _, v := range values {
		errors = append(errors, validationError{Field: v.Key, Type: v.Value})
	}

	respond(ctx, w, httpCode(ctx, ValidationError), errorResponse[[]validationError]{
		Type:    ValidationError,
		Code:    code,
		Message: "validation error",
		Data:    errors,
	})
}

func RespondWithJSON[T any](ctx context.Context, w http.ResponseWriter, payload T) {
	respond(ctx, w, http.StatusOK, payload)
}

func respond[T any](ctx context.Context, w http.ResponseWriter, status int, payload T) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Timestamp", time.Now().String())
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.ErrorContext(ctx, "failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
