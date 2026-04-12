package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
)

type TransactionService interface {
	List(context.Context, *uuid.UUID, int) ([]*model.Transaction, bool, error)
	Create(context.Context, *model.TransactionCreatePayload) (*model.Transaction, error)
}

type transactionService struct {
	transactions []*model.Transaction
}

func NewTransactionService(transactions []*model.Transaction) TransactionService {
	return &transactionService{transactions: transactions}
}

// Create implements [TransactionService].
func (t *transactionService) Create(ctx context.Context, payload *model.TransactionCreatePayload) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:     uuid.New(),
		Amount: payload.Amount,
	}
	t.transactions = append(t.transactions, transaction)
	return transaction, nil
}

// List implements [TransactionService].
func (t *transactionService) List(ctx context.Context, cursor *uuid.UUID, limit int) ([]*model.Transaction, bool, error) {
	start := 0
	if cursor == nil {
		start = 0
	} else {
		found := false
		for i, transaction := range t.transactions {
			if transaction.ID == *cursor {
				start = i + 1
				found = true
				break
			}
		}
		if !found {
			return nil, false, fmt.Errorf("cursor not found")
		}
	}

	result := t.transactions[start:min(start+limit, len(t.transactions))]
	return result, len(t.transactions) > start+limit, nil
}
