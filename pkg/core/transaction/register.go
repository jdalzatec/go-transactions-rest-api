package transaction

import (
	"database/sql"
	"net/http"

	"github.com/jdalzatec/banking/pkg/core/transaction/internal/handler"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/repository"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/service"
)

func NewServeMux(db *sql.DB) *http.ServeMux {
	repository := repository.NewSQLRepository(db)
	transactionService := service.NewTransactionService(repository)
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
