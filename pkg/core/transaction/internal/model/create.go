package model

type TransactionCreatePayload struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}
