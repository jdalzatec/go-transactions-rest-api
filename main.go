package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jdalzatec/banking/pkg/common/rest"
)

type Transaction struct {
	ID     uuid.UUID `json:"id"`
	Amount float64   `json:"amount"`
}

type TransactionCreatePayload struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

func main() {
	transactions := make([]Transaction, 0, 2)
	http.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		rest.JSON(r.Context(), w, http.StatusOK, transactions)
	})

	http.HandleFunc("POST /transactions", func(w http.ResponseWriter, r *http.Request) {
		var payload TransactionCreatePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			slog.Error("error decoding transaction create payload", "error", err)
			rest.Error(r.Context(), w, http.StatusBadRequest, rest.WithDetails("invalid payload"))
			return
		}
		if err := validator.New().Struct(payload); err != nil {
			rest.Error(r.Context(), w, http.StatusBadRequest, rest.WithDetails(err.Error()))
			return
		}

		data := Transaction{ID: uuid.New(), Amount: payload.Amount}
		transactions = append(transactions, data)
		rest.JSON(r.Context(), w, http.StatusCreated, data)
	})

	slog.Info("server starting", "port", ":1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		slog.Error("error running server", "error", err)
	}
}
