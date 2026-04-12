package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jdalzatec/banking/pkg/common/rest"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
)

type TransactionCreateHandler struct {
	Transactions []model.Transaction
}

func (h *TransactionCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload model.TransactionCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("error decoding transaction create payload", "error", err)
		rest.Error(r.Context(), w, http.StatusBadRequest, rest.WithDetails("invalid payload"))
		return
	}
	if err := validator.New().Struct(payload); err != nil {
		rest.Error(r.Context(), w, http.StatusBadRequest, rest.WithDetails(err.Error()))
		return
	}

	data := model.Transaction{ID: uuid.New(), Amount: payload.Amount}
	h.Transactions = append(h.Transactions, data)
	rest.JSON(r.Context(), w, http.StatusCreated, data)
}
