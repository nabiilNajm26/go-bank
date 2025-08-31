package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/repository"
	"github.com/shopspring/decimal"
)

var (
	ErrSameAccount = errors.New("cannot transfer to same account")
	ErrInvalidAmount = errors.New("invalid transfer amount")
)

type TransactionUseCase struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	db              *sqlx.DB
}

func NewTransactionUseCase(transactionRepo repository.TransactionRepository, accountRepo repository.AccountRepository, db *sqlx.DB) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		db:              db,
	}
}

func (uc *TransactionUseCase) Transfer(ctx context.Context, req *domain.TransferRequest) (*domain.Transaction, error) {
	fromAccountID, _ := uuid.Parse(req.FromAccountID)
	toAccountID, _ := uuid.Parse(req.ToAccountID)

	if fromAccountID == toAccountID {
		return nil, ErrSameAccount
	}

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}

	tx, err := uc.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Lock and get from account
	var fromAccount domain.Account
	err = tx.Get(&fromAccount, "SELECT * FROM accounts WHERE id = $1 FOR UPDATE", fromAccountID)
	if err != nil {
		return nil, err
	}

	// Lock and get to account
	var toAccount domain.Account
	err = tx.Get(&toAccount, "SELECT * FROM accounts WHERE id = $1 FOR UPDATE", toAccountID)
	if err != nil {
		return nil, err
	}

	// Check balance
	if fromAccount.Balance.LessThan(req.Amount) {
		return nil, ErrInsufficientBalance
	}

	// Update balances
	fromAccount.Balance = fromAccount.Balance.Sub(req.Amount)
	toAccount.Balance = toAccount.Balance.Add(req.Amount)

	_, err = tx.Exec("UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		fromAccount.Balance, fromAccount.ID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		toAccount.Balance, toAccount.ID)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	transaction := &domain.Transaction{
		ID:            uuid.New(),
		FromAccountID: &fromAccountID,
		ToAccountID:   &toAccountID,
		Amount:        req.Amount,
		Currency:      fromAccount.Currency,
		Type:          domain.TransactionTypeTransfer,
		Status:        domain.TransactionStatusCompleted,
		Reference:     uc.generateReference(),
		Description:   &req.Description,
		CreatedAt:     time.Now(),
	}

	completedAt := time.Now()
	transaction.CompletedAt = &completedAt

	_, err = tx.Exec(`
		INSERT INTO transactions (id, from_account_id, to_account_id, amount, currency, type, status, reference, description, created_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		transaction.ID, transaction.FromAccountID, transaction.ToAccountID, transaction.Amount,
		transaction.Currency, transaction.Type, transaction.Status, transaction.Reference,
		transaction.Description, transaction.CreatedAt, transaction.CompletedAt)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (uc *TransactionUseCase) GetTransactionHistory(ctx context.Context, accountID uuid.UUID, filter *domain.TransactionFilter) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetByAccountID(ctx, accountID, filter)
}

func (uc *TransactionUseCase) generateReference() string {
	return fmt.Sprintf("TXN%d", time.Now().Unix())
}