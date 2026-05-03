package model

import "github.com/oklog/ulid/v2"

type Transaction struct {
	ID     ulid.ULID `json:"id"`
	Amount float64   `json:"amount"`
}
