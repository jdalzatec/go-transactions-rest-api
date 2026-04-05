package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Transaction struct {
	ID     uuid.UUID `json:"id"`
	Amount float64   `json:"amount"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func responseWithError(w http.ResponseWriter, status int) {
	bytes, err := json.Marshal(ErrorResponse{Code: status, Message: http.StatusText(status)})
	if err != nil {
		slog.Error("error serializing", "error", err)
		http.Error(w, `{"code": 500, "message": "internal server error}`, http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(bytes))
}

func respondWithJSON(w http.ResponseWriter, status int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("error serializing", "error", err)
		responseWithError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(bytes))
}

func main() {
	http.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		data := []Transaction{
			{ID: uuid.New(), Amount: 100.0},
			{ID: uuid.New(), Amount: 200.0},
		}
		respondWithJSON(w, http.StatusOK, data)
	})

	http.HandleFunc("POST /transactions", func(w http.ResponseWriter, r *http.Request) {
		data := Transaction{ID: uuid.New(), Amount: 0.0}
		respondWithJSON(w, http.StatusCreated, data)
	})

	slog.Info("server starting", "port", ":1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		slog.Error("error running server", "error", err)
	}
}
