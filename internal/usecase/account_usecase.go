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
	ErrAccountNotEmpty = errors.New("account has non-zero balance")
	ErrUnauthorized = errors.New("unauthorized access")
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

func (uc *AccountUseCase) UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, req *domain.UpdateAccountRequest) (*domain.Account, error) {
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	// Check if user owns this account
	if account.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Update fields if provided
	if req.AccountType != nil {
		account.AccountType = *req.AccountType
	}
	if req.Status != nil {
		account.Status = *req.Status
	}
	account.UpdatedAt = time.Now()

	if err := uc.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (uc *AccountUseCase) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return ErrAccountNotFound
	}

	// Check if user owns this account
	if account.UserID != userID {
		return ErrUnauthorized
	}

	// Safety check: don't delete accounts with balance
	if !account.Balance.IsZero() {
		return ErrAccountNotEmpty
	}

	return uc.accountRepo.Delete(ctx, accountID)
}

func (uc *AccountUseCase) generateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%010d", rand.Intn(10000000000))
}