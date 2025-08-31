package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/infrastructure/session"
	"github.com/nabiilNajm26/go-bank/internal/repository"
	"github.com/nabiilNajm26/go-bank/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthUseCase struct {
	userRepo       repository.UserRepository
	jwtManager     *utils.JWTManager
	sessionService *session.SessionService
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtManager *utils.JWTManager, sessionService *session.SessionService) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		jwtManager:     jwtManager,
		sessionService: sessionService,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.AuthResponse, error) {
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        &req.Phone,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	var accessToken, refreshToken string
	
	if uc.sessionService != nil {
		sessionID, err := uc.sessionService.CreateSession(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		accessToken, err = uc.sessionService.GenerateTokenWithSession(user.ID, sessionID, uc.jwtManager.GetAccessSecret(), time.Hour)
		if err != nil {
			return nil, err
		}

		refreshToken, err = uc.sessionService.GenerateTokenWithSession(user.ID, sessionID, uc.jwtManager.GetRefreshSecret(), 7*24*time.Hour)
		if err != nil {
			return nil, err
		}
	} else {
		accessToken, refreshToken, err = uc.jwtManager.GenerateTokenPair(user.ID, user.Email)
		if err != nil {
			return nil, err
		}
	}

	return &domain.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	var accessToken, refreshToken string
	
	if uc.sessionService != nil {
		sessionID, err := uc.sessionService.CreateSession(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		accessToken, err = uc.sessionService.GenerateTokenWithSession(user.ID, sessionID, uc.jwtManager.GetAccessSecret(), time.Hour)
		if err != nil {
			return nil, err
		}

		refreshToken, err = uc.sessionService.GenerateTokenWithSession(user.ID, sessionID, uc.jwtManager.GetRefreshSecret(), 7*24*time.Hour)
		if err != nil {
			return nil, err
		}
	} else {
		accessToken, refreshToken, err = uc.jwtManager.GenerateTokenPair(user.ID, user.Email)
		if err != nil {
			return nil, err
		}
	}

	return &domain.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
	}, nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	claims, err := uc.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	newAccessToken, newRefreshToken, err := uc.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:         user,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600,
	}, nil
}