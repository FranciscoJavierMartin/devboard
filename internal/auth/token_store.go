package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrTokenRevoked when token is revoked
var ErrTokenRevoked = errors.New("revoked token")

// ErrTokenNotFound when token is not found (on Redis instance)
var ErrTokenNotFound = errors.New("refresh token not found")

type TokenStore struct {
	redis *redis.Client
}

func NewTokenStore(redisClient *redis.Client) *TokenStore {
	return &TokenStore{redis: redisClient}
}

func refreshKey(tokenID string) string {
	return "refresh:" + tokenID
}

func revokedKey(tokenID string) string {
	return "revoked:" + tokenID
}

func (s *TokenStore) StoreRefreshToken(ctx context.Context, tokenID, userID string, ttl time.Duration) error {
	err := s.redis.Set(
		ctx,
		refreshKey(tokenID),
		userID,
		ttl,
	).Err()

	if err != nil {
		return fmt.Errorf("saving refresh token on Redis instance: %w", err)
	}

	return nil
}

func (s *TokenStore) ValidateRefreshToken(ctx context.Context, tokenID string) (string, error) {
	revoked, err := s.redis.Exists(ctx, revokedKey(tokenID)).Result()

	if err != nil {
		return "", fmt.Errorf("verifying revoking: %w", err)
	}

	if revoked > 0 {
		return "", ErrTokenRevoked
	}

	userID, err := s.redis.Get(ctx, refreshKey(tokenID)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenNotFound
		}

		return "", fmt.Errorf("retrieving refresh token: %w", err)
	}

	return userID, nil
}

func (s *TokenStore) RevokeRefreshToken(ctx context.Context, tokenID string, ttl time.Duration) error {
	pipe := s.redis.Pipeline()

	pipe.Set(ctx, revokedKey(tokenID), "1", ttl)
	pipe.Del(ctx, refreshKey(tokenID))

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("revoking refresh token: %w", err)
	}

	return nil
}

func (s *TokenStore) RevokeAllUserTokens(ctx context.Context, userID string, ttl time.Duration) error {
	var cursor uint64
	pattern := "refresh:*"

	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("searching user tokens: %w", err)
		}

		for _, key := range keys {
			storedUserID, err := s.redis.Get(ctx, key).Result()

			if err != nil {
				continue
			}

			if storedUserID == userID {
				tokenID := key[len("refresh:"):]
				_ = s.RevokeRefreshToken(ctx, tokenID, ttl)
			}
		}

		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return nil
}
