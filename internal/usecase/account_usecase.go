package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/repository"
	"github.com/shopspring/decimal"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type AccountUseCase struct {
	accountRepo repository.AccountRepository
	userRepo    repository.UserRepository
}

func NewAccountUseCase(accountRepo repository.AccountRepository, userRepo repository.UserRepository) *AccountUseCase {
	return &AccountUseCase{
		accountRepo: accountRepo,
		userRepo:    userRepo,
	}
}

func (uc *AccountUseCase) CreateAccount(ctx context.Context, userID uuid.UUID, req *domain.CreateAccountRequest) (*domain.Account, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	account := &domain.Account{
		ID:            uuid.New(),
		UserID:        userID,
		AccountNumber: uc.generateAccountNumber(),
		AccountType:   req.AccountType,
		Balance:       decimal.NewFromInt(0),
		Currency:      req.Currency,
		Status:        domain.AccountStatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (uc *AccountUseCase) GetAccount(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (uc *AccountUseCase) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	return uc.accountRepo.GetByUserID(ctx, userID)
}

func (uc *AccountUseCase) generateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%010d", rand.Intn(10000000000))
}