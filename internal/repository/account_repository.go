package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}