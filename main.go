package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Transaction struct {
	ID     uuid.UUID `json:"id"`
	Amount float64   `json:"amount"`
}

type TransactionCreatePayload struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type ErrorOption func(*ErrorResponse)

func withDetails(details string) ErrorOption {
	return func(er *ErrorResponse) {
		er.Details = details
	}
}

func responseWithError(w http.ResponseWriter, status int, options ...ErrorOption) {
	errorResponse := ErrorResponse{Code: status, Message: http.StatusText(status)}
	for _, option := range options {
		option(&errorResponse)
	}

	bytes, err := json.Marshal(errorResponse)
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
	transactions := make([]Transaction, 0, 2)
	http.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, transactions)
	})

	http.HandleFunc("POST /transactions", func(w http.ResponseWriter, r *http.Request) {
		var payload TransactionCreatePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			slog.Error("error decoding transaction create payload", "error", err)
			responseWithError(w, http.StatusBadRequest, withDetails("invalid payload"))
			return
		}
		if err := validator.New().Struct(payload); err != nil {
			responseWithError(
				w,
				http.StatusBadRequest,
				withDetails(err.Error()),
			)
			return
		}

		data := Transaction{ID: uuid.New(), Amount: payload.Amount}
		transactions = append(transactions, data)
		respondWithJSON(w, http.StatusCreated, data)
	})

	slog.Info("server starting", "port", ":1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		slog.Error("error running server", "error", err)
	}
}
