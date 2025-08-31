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

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) repository.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `
		INSERT INTO accounts (user_id, account_number, account_type, balance, currency, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		account.UserID,
		account.AccountNumber,
		account.AccountType,
		account.Balance,
		account.Currency,
		account.Status,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	return err
}

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var account domain.Account
	query := `SELECT * FROM accounts WHERE id = $1`
	
	err := r.db.GetContext(ctx, &account, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	var account domain.Account
	query := `SELECT * FROM accounts WHERE account_number = $1`
	
	err := r.db.GetContext(ctx, &account, query, accountNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	var accounts []*domain.Account
	query := `SELECT * FROM accounts WHERE user_id = $1 ORDER BY created_at DESC`
	
	err := r.db.SelectContext(ctx, &accounts, query, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *accountRepository) Update(ctx context.Context, account *domain.Account) error {
	query := `
		UPDATE accounts 
		SET account_type = $2, balance = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.AccountType,
		account.Balance,
		account.Status,
	)

	return err
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}