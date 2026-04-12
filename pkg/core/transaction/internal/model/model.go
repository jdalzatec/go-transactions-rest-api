package model

import "github.com/google/uuid"

type Transaction struct {
	ID     uuid.UUID `json:"id"`
	Amount float64   `json:"amount"`
}
