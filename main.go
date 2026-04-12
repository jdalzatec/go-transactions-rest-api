package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jdalzatec/banking/pkg/common/rest"
	"github.com/jdalzatec/banking/pkg/core/transaction"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/transactions", transaction.NewServeMux())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rest.Error(
			r.Context(),
			w,
			http.StatusNotFound,
			rest.WithDetails(fmt.Sprintf("%s %s not found", r.Method, r.URL.Path)),
		)
	})

	slog.Info("server starting", "port", ":1234")
	err := http.ListenAndServe(":1234", mux)
	if err != nil {
		slog.Error("error running server", "error", err)
	}
}
