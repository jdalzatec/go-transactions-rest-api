package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type ErrorOption func(*ErrorResponse)

func WithDetails(details string) ErrorOption {
	return func(er *ErrorResponse) {
		er.Details = details
	}
}

func Error(ctx context.Context, w http.ResponseWriter, status int, options ...ErrorOption) {
	errorResponse := ErrorResponse{Code: status, Message: strings.ToLower(http.StatusText(status))}
	for _, option := range options {
		option(&errorResponse)
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		slog.ErrorContext(ctx, "error serializing", "error", err)
		h := w.Header()
		h.Del("Content-Length")
		h.Set("Content-Type", "application/json")
		h.Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintln(w, `{"code": 500, "message": "internal server error"}`)
		if err != nil {
			slog.ErrorContext(ctx, "error writing response", "error", err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(bytes)
	if err != nil {
		slog.ErrorContext(ctx, "error writing response", "error", err)
	}
}
