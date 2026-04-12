package transaction

import (
	"net/http"

	"github.com/jdalzatec/banking/pkg/core/transaction/internal/handler"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
)

func NewServeMux() *http.ServeMux {
	transactions := make([]model.Transaction, 0, 2)
	mux := http.NewServeMux()
	mux.Handle(
		"GET /",
		&handler.TransactionListHandler{Transactions: transactions},
	)
	mux.Handle(
		"POST /",
		&handler.TransactionCreateHandler{Transactions: transactions},
	)
	return mux
}
