package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jdalzatec/banking/pkg/common/rest"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/service"
)

type TransactionCreateHandler struct {
	TransactionService service.TransactionService
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

	transaction, err := h.TransactionService.Create(r.Context(), &payload)
	if err != nil {
		rest.Error(r.Context(), w, http.StatusInternalServerError, rest.WithDetails(err.Error()))
		return
	}
	rest.JSON(r.Context(), w, http.StatusCreated, transaction)
}
