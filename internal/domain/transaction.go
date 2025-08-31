package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypePayment    TransactionType = "payment"

	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

type Transaction struct {
	ID            uuid.UUID         `json:"id" db:"id"`
	FromAccountID *uuid.UUID        `json:"from_account_id,omitempty" db:"from_account_id"`
	ToAccountID   *uuid.UUID        `json:"to_account_id,omitempty" db:"to_account_id"`
	Amount        decimal.Decimal   `json:"amount" db:"amount"`
	Currency      string            `json:"currency" db:"currency"`
	Type          TransactionType   `json:"type" db:"type"`
	Status        TransactionStatus `json:"status" db:"status"`
	Reference     string            `json:"reference" db:"reference"`
	Description   *string           `json:"description,omitempty" db:"description"`
	Metadata      map[string]any    `json:"metadata,omitempty" db:"metadata"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	CompletedAt   *time.Time        `json:"completed_at,omitempty" db:"completed_at"`
}

type TransferRequest struct {
	FromAccountID string          `json:"from_account_id" validate:"required,uuid"`
	ToAccountID   string          `json:"to_account_id" validate:"required,uuid"`
	Amount        decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Description   string          `json:"description,omitempty" validate:"omitempty,max=500"`
}

type TransactionFilter struct {
	AccountID uuid.UUID
	Type      TransactionType
	Status    TransactionStatus
	FromDate  time.Time
	ToDate    time.Time
	Limit     int
	Offset    int
}