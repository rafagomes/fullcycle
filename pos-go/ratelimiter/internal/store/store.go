package store

import (
	"context"
	"time"
)

type RateLimiterStore interface {
	Increment(ctx context.Context, key string, expiration time.Duration) (int64, error)

	Block(ctx context.Context, identifier string, duration time.Duration) error

	IsBlocked(ctx context.Context, identifier string) (bool, error)
}
