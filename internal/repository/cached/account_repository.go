package cached

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/infrastructure/cache"
	"github.com/nabiilNajm26/go-bank/internal/repository"
)

type cachedAccountRepository struct {
	repo  repository.AccountRepository
	cache *cache.CacheService
}

func NewCachedAccountRepository(repo repository.AccountRepository, cache *cache.CacheService) repository.AccountRepository {
	return &cachedAccountRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *cachedAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	err := r.repo.Create(ctx, account)
	if err != nil {
		return err
	}

	// Cache the created account
	if err := r.cache.SetAccount(ctx, account); err != nil {
		log.Printf("Failed to cache account after creation: %v", err)
	}

	return nil
}

func (r *cachedAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	// Try cache first
	account, err := r.cache.GetAccount(ctx, id)
	if err == nil && account != nil {
		return account, nil
	}

	// If not in cache or error (except Redis not available), fetch from DB
	if err != nil && err != redis.Nil {
		log.Printf("Cache error for account %s: %v", id, err)
	}

	account, err = r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if account != nil {
		// Cache the result
		if err := r.cache.SetAccount(ctx, account); err != nil {
			log.Printf("Failed to cache account: %v", err)
		}
	}

	return account, nil
}

func (r *cachedAccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	// For user's accounts list, we don't cache (could change frequently)
	return r.repo.GetByUserID(ctx, userID)
}

func (r *cachedAccountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	// Account number lookups bypass cache for now
	account, err := r.repo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, err
	}

	if account != nil {
		// Cache the result
		if err := r.cache.SetAccount(ctx, account); err != nil {
			log.Printf("Failed to cache account after account number lookup: %v", err)
		}
	}

	return account, nil
}

func (r *cachedAccountRepository) Update(ctx context.Context, account *domain.Account) error {
	err := r.repo.Update(ctx, account)
	if err != nil {
		return err
	}

	// Invalidate cache
	if err := r.cache.DeleteAccount(ctx, account.ID); err != nil {
		log.Printf("Failed to invalidate account cache: %v", err)
	}

	return nil
}

func (r *cachedAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache
	if err := r.cache.DeleteAccount(ctx, id); err != nil {
		log.Printf("Failed to invalidate account cache: %v", err)
	}

	return nil
}

