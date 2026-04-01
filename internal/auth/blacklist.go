package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlacklist struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewTokenBlacklist(rdb *redis.Client, ttl time.Duration) *TokenBlacklist {
	return &TokenBlacklist{rdb: rdb, ttl: ttl}
}

func (b *TokenBlacklist) Add(ctx context.Context, jti string, exp time.Time) error {
	key := fmt.Sprintf("jwt:blacklist:%s", jti)
	ttl := time.Until(exp)
	if ttl <= 0 {
		return nil
	}
	if b.ttl > 0 && ttl > b.ttl {
		ttl = b.ttl
	}
	return b.rdb.Set(ctx, key, "1", ttl).Err()
}

func (b *TokenBlacklist) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("jwt:blacklist:%s", jti)
	n, err := b.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
