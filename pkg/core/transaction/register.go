package transaction

import (
	"net/http"

	"github.com/jdalzatec/banking/pkg/core/transaction/internal/handler"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/service"
)

func NewServeMux() *http.ServeMux {
	transactions := make([]*model.Transaction, 0)
	transactionService := service.NewTransactionService(transactions)
	mux := http.NewServeMux()
	mux.Handle(
		"GET /",
		&handler.TransactionListHandler{TransactionService: transactionService},
	)
	mux.Handle(
		"POST /",
		&handler.TransactionCreateHandler{TransactionService: transactionService},
	)
	return mux
}
