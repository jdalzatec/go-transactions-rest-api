package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jdalzatec/banking/pkg/common/rest"
	"github.com/jdalzatec/banking/pkg/core/transaction"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/transactions", transaction.NewServeMux(db))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rest.Error(
			r.Context(),
			w,
			http.StatusNotFound,
			rest.WithDetails(fmt.Sprintf("%s %s not found", r.Method, r.URL.Path)),
		)
	})

	slog.Info("server starting", "port", ":1234")
	if err := http.ListenAndServe(":1234", mux); err != nil {
		slog.Error("error running server", "error", err)
	}
}
