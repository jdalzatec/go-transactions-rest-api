package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jdalzatec/banking/pkg/common/pagination"
	"github.com/jdalzatec/banking/pkg/common/rest"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/service"
	"github.com/oklog/ulid/v2"
)

type TransactionListHandler struct {
	TransactionService service.TransactionService
}

func (h *TransactionListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var cursor *ulid.ULID
	cursorValue := r.URL.Query().Get("cursor")
	if cursorValue == "" {
		cursor = nil
	} else {
		parsedCursor, cursorErr := ulid.Parse(cursorValue)
		if cursorErr != nil {
			rest.Error(r.Context(), w, http.StatusBadRequest, rest.WithDetails(fmt.Sprintf("invalid cursor: %s", cursorErr.Error())))
			return
		} else {
			cursor = &parsedCursor
		}
	}

	limit := 10
	var limitErr error
	limitValue := r.URL.Query().Get("limit")
	if limitValue != "" {
		limit, limitErr = strconv.Atoi(limitValue)
		if limitErr != nil {
			rest.Error(r.Context(), w, http.StatusBadRequest, rest.WithDetails(fmt.Sprintf("invalid limit: %s", limitErr.Error())))
			return
		}
		if limit <= 0 {
			rest.Error(r.Context(), w, http.StatusBadRequest, rest.WithDetails("limit must be greater than 0"))
			return
		}
	}

	transactions, hasMore, err := h.TransactionService.List(r.Context(), cursor, limit)
	if err != nil {
		rest.Error(r.Context(), w, http.StatusInternalServerError, rest.WithDetails(err.Error()))
		return
	}

	if len(transactions) == 0 {
		rest.JSON(r.Context(), w, http.StatusOK, pagination.Paginated[model.Transaction]{
			Items:      []*model.Transaction{},
			Pagination: pagination.Pagination{Cursor: nil, HasMore: hasMore},
		})
		return
	}

	paginated := pagination.Paginated[model.Transaction]{
		Items: transactions,
		Pagination: pagination.Pagination{
			Cursor:  &transactions[len(transactions)-1].ID,
			HasMore: hasMore,
		},
	}
	rest.JSON(r.Context(), w, http.StatusOK, paginated)
}
