package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/repository"
)

var (
	ErrUserHasActiveAccounts = errors.New("user has active accounts")
)

type UserUseCase struct {
	userRepo    repository.UserRepository
	accountRepo repository.AccountRepository
}

func NewUserUseCase(userRepo repository.UserRepository, accountRepo repository.AccountRepository) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		accountRepo: accountRepo,
	}
}

func (uc *UserUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUseCase) UpdateProfileImage(ctx context.Context, userID uuid.UUID, imageURL string) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.ProfileImageURL = &imageURL
	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, userID uuid.UUID, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check email uniqueness if email is being updated
	if req.Email != nil && *req.Email != user.Email {
		existing, err := uc.userRepo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrEmailAlreadyExists
		}
		user.Email = *req.Email
	}

	// Update other fields if provided
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Check if user has active accounts
	if uc.accountRepo != nil {
		accounts, err := uc.accountRepo.GetByUserID(ctx, userID)
		if err != nil {
			return err
		}
		if len(accounts) > 0 {
			return ErrUserHasActiveAccounts
		}
	}

	return uc.userRepo.Delete(ctx, userID)
}