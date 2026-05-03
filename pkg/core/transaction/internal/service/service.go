package service

import (
	"context"
	"fmt"

	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
	"github.com/jdalzatec/banking/pkg/core/transaction/internal/repository"
	"github.com/oklog/ulid/v2"
)

type TransactionService interface {
	List(context.Context, *ulid.ULID, int) ([]*model.Transaction, bool, error)
	Create(context.Context, *model.TransactionCreatePayload) (*model.Transaction, error)
}

type transactionService struct {
	repository repository.Repository
}

func NewTransactionService(repository repository.Repository) TransactionService {
	return &transactionService{repository: repository}
}

func (service *transactionService) Create(ctx context.Context, payload *model.TransactionCreatePayload) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:     ulid.Make(),
		Amount: payload.Amount,
	}
	outout, err := service.repository.Save(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction from service: %w", err)
	}

	return outout, nil
}

func (service *transactionService) List(ctx context.Context, cursor *ulid.ULID, limit int) ([]*model.Transaction, bool, error) {
	var concreteCursor ulid.ULID
	if cursor != nil {
		concreteCursor = *cursor
	}
	result, err := service.repository.FindMany(ctx, concreteCursor, limit)
	if err != nil {
		return nil, true, fmt.Errorf("error listing transactions from service: %w", err)
	}
	return result, true, nil
}
