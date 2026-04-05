package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Transaction struct {
	ID     uuid.UUID `json:"id"`
	Amount float64   `json:"amount"`
}

func main() {
	http.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		data := []Transaction{
			{ID: uuid.New(), Amount: 100.0},
			{ID: uuid.New(), Amount: 200.0},
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			// TODO: ofuscate the error before logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(bytes))
	})
	slog.Info("server starting", "port", ":1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		slog.Error("error running server", "error", err)
	}
}
