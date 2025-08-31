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

type cachedUserRepository struct {
	repo  repository.UserRepository
	cache *cache.CacheService
}

func NewCachedUserRepository(repo repository.UserRepository, cache *cache.CacheService) repository.UserRepository {
	return &cachedUserRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *cachedUserRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	// Cache the created user
	if err := r.cache.SetUser(ctx, user); err != nil {
		log.Printf("Failed to cache user after creation: %v", err)
	}

	return nil
}

func (r *cachedUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	// Try cache first
	user, err := r.cache.GetUser(ctx, id)
	if err == nil && user != nil {
		return user, nil
	}

	// If not in cache or error (except Redis not available), fetch from DB
	if err != nil && err != redis.Nil {
		log.Printf("Cache error for user %s: %v", id, err)
	}

	user, err = r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		// Cache the result
		if err := r.cache.SetUser(ctx, user); err != nil {
			log.Printf("Failed to cache user: %v", err)
		}
	}

	return user, nil
}

func (r *cachedUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Email lookups bypass cache for now (could implement email-based cache key)
	user, err := r.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		// Cache the result
		if err := r.cache.SetUser(ctx, user); err != nil {
			log.Printf("Failed to cache user after email lookup: %v", err)
		}
	}

	return user, nil
}

func (r *cachedUserRepository) Update(ctx context.Context, user *domain.User) error {
	err := r.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Invalidate cache
	if err := r.cache.DeleteUser(ctx, user.ID); err != nil {
		log.Printf("Failed to invalidate user cache: %v", err)
	}

	return nil
}

func (r *cachedUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache
	if err := r.cache.DeleteUser(ctx, id); err != nil {
		log.Printf("Failed to invalidate user cache: %v", err)
	}

	return nil
}