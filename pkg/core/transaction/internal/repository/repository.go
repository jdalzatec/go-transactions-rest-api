package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jdalzatec/banking/pkg/core/transaction/internal/model"
	"github.com/oklog/ulid/v2"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type Repository interface {
	FindMany(ctx context.Context, cursor ulid.ULID, limit int) ([]*model.Transaction, error)
	Find(ctx context.Context, id ulid.ULID) (*model.Transaction, error)
	Save(ctx context.Context, t *model.Transaction) (*model.Transaction, error)
}

type SQLRepository struct {
	db *sql.DB
}

func (r *SQLRepository) FindMany(ctx context.Context, cursor ulid.ULID, limit int) ([]*model.Transaction, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, amount FROM transactions WHERE id > $1 LIMIT $2", cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying the database: %w", err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error("error closing session", "error", err)
		}
	}()

	results := make([]*model.Transaction, 0)
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.Amount); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		results = append(results, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (r *SQLRepository) Find(ctx context.Context, id ulid.ULID) (*model.Transaction, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, amount FROM transactions WHERE id = $1", id)
	var t model.Transaction
	if err := row.Scan(&id, &t.Amount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("error scanning row: %w", err)
	}

	return &t, nil
}

func (r *SQLRepository) Save(ctx context.Context, t *model.Transaction) (*model.Transaction, error) {
	row := r.db.QueryRowContext(ctx, "INSERT INTO transactions (id, amount) VALUES ($1, $2) RETURNING id, amount", t.ID, t.Amount)
	var output model.Transaction
	if err := row.Scan(&output.ID, &output.Amount); err != nil {
		return nil, fmt.Errorf("error creating row: %w", err)
	}
	return &output, nil
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}
