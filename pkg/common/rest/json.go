package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

func JSON(ctx context.Context, w http.ResponseWriter, status int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		slog.ErrorContext(ctx, "error serializing", "error", err)
		Error(ctx, w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(bytes)
	if err != nil {
		slog.ErrorContext(ctx, "error writing response", "error", err)
	}
}
