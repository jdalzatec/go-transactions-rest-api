package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jdalzatec/banking/pkg/common/pagination"
	"github.com/jdalzatec/banking/pkg/common/rest"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
)

type TransactionListHandler struct {
	Transactions []model.Transaction
}

func (h *TransactionListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(h.Transactions) == 0 {
		rest.JSON(r.Context(), w, http.StatusOK, pagination.Paginated[model.Transaction]{
			Items:      []model.Transaction{},
			Pagination: pagination.Pagination{Cursor: uuid.Nil, HasMore: false},
		})
		return
	}

	paginated := pagination.Paginated[model.Transaction]{
		Items: h.Transactions,
		Pagination: pagination.Pagination{
			Cursor:  h.Transactions[len(h.Transactions)-1].ID,
			HasMore: true,
		},
	}
	rest.JSON(r.Context(), w, http.StatusOK, paginated)
}
