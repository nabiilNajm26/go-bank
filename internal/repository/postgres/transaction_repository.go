package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/repository"
)

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) repository.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (from_account_id, to_account_id, amount, currency, type, status, reference, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		tx.FromAccountID,
		tx.ToAccountID,
		tx.Amount,
		tx.Currency,
		tx.Type,
		tx.Status,
		tx.Reference,
		tx.Description,
	).Scan(&tx.ID, &tx.CreatedAt)

	return err
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var tx domain.Transaction
	query := `SELECT * FROM transactions WHERE id = $1`
	
	err := r.db.GetContext(ctx, &tx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &tx, nil
}

func (r *transactionRepository) GetByReference(ctx context.Context, reference string) (*domain.Transaction, error) {
	var tx domain.Transaction
	query := `SELECT * FROM transactions WHERE reference = $1`
	
	err := r.db.GetContext(ctx, &tx, query, reference)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &tx, nil
}

func (r *transactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID, filter *domain.TransactionFilter) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	query := `
		SELECT * FROM transactions 
		WHERE (from_account_id = $1 OR to_account_id = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	limit := 50
	offset := 0
	if filter != nil {
		if filter.Limit > 0 {
			limit = filter.Limit
		}
		offset = filter.Offset
	}

	err := r.db.SelectContext(ctx, &transactions, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	query := `
		UPDATE transactions 
		SET status = $2, completed_at = $3
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, tx.ID, tx.Status, tx.CompletedAt)
	return err
}