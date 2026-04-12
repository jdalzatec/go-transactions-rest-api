package transaction

import (
	"net/http"

	"github.com/jdalzatec/banking/pkg/core/transaction/internal/handler"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/service"
)

func NewServeMux() *http.ServeMux {
	transactions := make([]*model.Transaction, 0, 2)
	mux := http.NewServeMux()
	mux.Handle(
		"GET /",
		&handler.TransactionListHandler{TransactionService: service.NewTransactionService(transactions)},
	)
	mux.Handle(
		"POST /",
		&handler.TransactionCreateHandler{TransactionService: service.NewTransactionService(transactions)},
	)
	return mux
}
