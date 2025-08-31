package session

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nabiilNajm26/go-bank/internal/infrastructure/cache"
)

type SessionService struct {
	cache *cache.CacheService
}

func NewSessionService(cache *cache.CacheService) *SessionService {
	return &SessionService{
		cache: cache,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	sessionID := uuid.New().String()
	
	err := s.cache.SetSession(ctx, sessionID, userID)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	return sessionID, nil
}

func (s *SessionService) GetUserFromSession(ctx context.Context, sessionID string) (uuid.UUID, error) {
	return s.cache.GetSession(ctx, sessionID)
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return s.cache.DeleteSession(ctx, sessionID)
}

func (s *SessionService) ExtractSessionFromToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sessionID, exists := claims["session_id"]; exists {
			return sessionID.(string), nil
		}
		return "", fmt.Errorf("session_id not found in token")
	}

	return "", fmt.Errorf("invalid token")
}

func (s *SessionService) GenerateTokenWithSession(userID uuid.UUID, sessionID string, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID.String(),
		"session_id": sessionID,
		"exp":        time.Now().Add(duration).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}