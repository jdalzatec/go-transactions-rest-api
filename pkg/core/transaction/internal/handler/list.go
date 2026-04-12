package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jdalzatec/banking/pkg/common/pagination"
	"github.com/jdalzatec/banking/pkg/common/rest"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/service"
)

type TransactionListHandler struct {
	TransactionService service.TransactionService
}

func (h *TransactionListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cursor, err := uuid.Parse(r.URL.Query().Get("cursor"))
	if err != nil {
		rest.Error(r.Context(), w, http.StatusBadRequest, rest.WithDetails(fmt.Sprintf("invalid cursor: %s", err.Error())))
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		rest.Error(r.Context(), w, http.StatusBadRequest, rest.WithDetails(fmt.Sprintf("invalid limit: %s", err.Error())))
		return
	}

	transactions, err := h.TransactionService.List(r.Context(), cursor, limit)
	if err != nil {
		rest.Error(r.Context(), w, http.StatusInternalServerError, rest.WithDetails(err.Error()))
		return
	}

	if len(transactions) == 0 {
		rest.JSON(r.Context(), w, http.StatusOK, pagination.Paginated[model.Transaction]{
			Items:      []*model.Transaction{},
			Pagination: pagination.Pagination{Cursor: uuid.Nil, HasMore: false},
		})
		return
	}

	paginated := pagination.Paginated[model.Transaction]{
		Items: transactions,
		Pagination: pagination.Pagination{
			Cursor:  transactions[len(transactions)-1].ID,
			HasMore: true,
		},
	}
	rest.JSON(r.Context(), w, http.StatusOK, paginated)
}
