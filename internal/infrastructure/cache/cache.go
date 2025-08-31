package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/domain"
	"github.com/nabiilNajm26/go-bank/internal/infrastructure/redis"
)

type CacheService struct {
	redis *redis.RedisClient
}

func NewCacheService(redisClient *redis.RedisClient) *CacheService {
	return &CacheService{
		redis: redisClient,
	}
}

// User caching
func (c *CacheService) SetUser(ctx context.Context, user *domain.User) error {
	key := fmt.Sprintf("user:%s", user.ID.String())
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, data, 30*time.Minute)
}

func (c *CacheService) GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	key := fmt.Sprintf("user:%s", userID.String())
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *CacheService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("user:%s", userID.String())
	return c.redis.Del(ctx, key)
}

// Account caching
func (c *CacheService) SetAccount(ctx context.Context, account *domain.Account) error {
	key := fmt.Sprintf("account:%s", account.ID.String())
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, data, 15*time.Minute)
}

func (c *CacheService) GetAccount(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	key := fmt.Sprintf("account:%s", accountID.String())
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var account domain.Account
	if err := json.Unmarshal([]byte(data), &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (c *CacheService) DeleteAccount(ctx context.Context, accountID uuid.UUID) error {
	key := fmt.Sprintf("account:%s", accountID.String())
	return c.redis.Del(ctx, key)
}

// Transaction caching for recent transactions
func (c *CacheService) SetTransactions(ctx context.Context, accountID uuid.UUID, transactions []*domain.Transaction) error {
	key := fmt.Sprintf("transactions:%s", accountID.String())
	data, err := json.Marshal(transactions)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, data, 5*time.Minute)
}

func (c *CacheService) GetTransactions(ctx context.Context, accountID uuid.UUID) ([]*domain.Transaction, error) {
	key := fmt.Sprintf("transactions:%s", accountID.String())
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var transactions []*domain.Transaction
	if err := json.Unmarshal([]byte(data), &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (c *CacheService) DeleteTransactions(ctx context.Context, accountID uuid.UUID) error {
	key := fmt.Sprintf("transactions:%s", accountID.String())
	return c.redis.Del(ctx, key)
}

// Session management
func (c *CacheService) SetSession(ctx context.Context, sessionID string, userID uuid.UUID) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.redis.Set(ctx, key, userID.String(), 24*time.Hour)
}

func (c *CacheService) GetSession(ctx context.Context, sessionID string) (uuid.UUID, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(data)
}

func (c *CacheService) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.redis.Del(ctx, key)
}

// Cache invalidation
func (c *CacheService) InvalidateUserCache(ctx context.Context, userID uuid.UUID) error {
	userKey := fmt.Sprintf("user:%s", userID.String())
	return c.redis.Del(ctx, userKey)
}

func (c *CacheService) InvalidateAccountCache(ctx context.Context, accountID uuid.UUID) error {
	accountKey := fmt.Sprintf("account:%s", accountID.String())
	transactionKey := fmt.Sprintf("transactions:%s", accountID.String())
	return c.redis.Del(ctx, accountKey, transactionKey)
}