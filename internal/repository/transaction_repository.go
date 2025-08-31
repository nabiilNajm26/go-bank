package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetByReference(ctx context.Context, reference string) (*domain.Transaction, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID, filter *domain.TransactionFilter) ([]*domain.Transaction, error)
	Update(ctx context.Context, tx *domain.Transaction) error
}