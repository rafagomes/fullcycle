package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string, db int) (RateLimiterStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisStore{client: client}, nil
}

func (r *RedisStore) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		r.client.Expire(ctx, key, expiration)
	}
	return count, nil
}

func (r *RedisStore) Block(ctx context.Context, identifier string, duration time.Duration) error {
	blockKey := "block:" + identifier
	return r.client.Set(ctx, blockKey, "1", duration).Err()
}

func (r *RedisStore) IsBlocked(ctx context.Context, identifier string) (bool, error) {
	blockKey := "block:" + identifier
	val, err := r.client.Get(ctx, blockKey).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == "1", nil
}
