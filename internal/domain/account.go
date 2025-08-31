package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountStatus string
type AccountType string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
	AccountStatusFrozen   AccountStatus = "frozen"
	AccountStatusClosed   AccountStatus = "closed"

	AccountTypeSavings AccountType = "savings"
	AccountTypeChecking AccountType = "checking"
	AccountTypeDeposit  AccountType = "deposit"
)

type Account struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	UserID        uuid.UUID       `json:"user_id" db:"user_id"`
	AccountNumber string          `json:"account_number" db:"account_number"`
	AccountType   AccountType     `json:"account_type" db:"account_type"`
	Balance       decimal.Decimal `json:"balance" db:"balance"`
	Currency      string          `json:"currency" db:"currency"`
	Status        AccountStatus   `json:"status" db:"status"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

type CreateAccountRequest struct {
	AccountType AccountType `json:"account_type" validate:"required,oneof=savings checking deposit"`
	Currency    string      `json:"currency" validate:"required,len=3"`
}

type AccountResponse struct {
	*Account
	User *User `json:"user,omitempty"`
}