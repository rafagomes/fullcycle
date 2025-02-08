package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter"
)

type FakeStore struct {
	mu       sync.Mutex
	counters map[string]int64
	blocks   map[string]time.Time
}

func NewFakeStore() *FakeStore {
	return &FakeStore{
		counters: make(map[string]int64),
		blocks:   make(map[string]time.Time),
	}
}

func (f *FakeStore) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.counters[key]++

	return f.counters[key], nil
}

func (f *FakeStore) Block(ctx context.Context, identifier string, duration time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.blocks[identifier] = time.Now().Add(duration)
	return nil
}

func (f *FakeStore) IsBlocked(ctx context.Context, identifier string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	t, exists := f.blocks[identifier]
	if !exists {
		return false, nil
	}
	if time.Now().After(t) {
		delete(f.blocks, identifier)
		return false, nil
	}
	return true, nil
}

func TestRateLimiter_IP(t *testing.T) {
	fakeStore := NewFakeStore()
	cfg := &config.Config{
		RateLimitIP:    5,
		RateLimitToken: 10,
		BlockTimeIP:    5 * time.Second,
		BlockTimeToken: 5 * time.Second,
	}
	rl := limiter.NewRateLimiter(fakeStore, cfg)
	ctx := context.Background()
	identifier := "ip:127.0.0.1"

	for i := 0; i < 5; i++ {
		allowed, err := rl.Allow(ctx, identifier, cfg.RateLimitIP, cfg.BlockTimeIP)
		if err != nil {
			t.Fatalf("Erro no Allow: %v", err)
		}
		if !allowed {
			t.Fatalf("Requisição %d foi bloqueada, mas deveria ser permitida", i+1)
		}
	}

	allowed, err := rl.Allow(ctx, identifier, cfg.RateLimitIP, cfg.BlockTimeIP)
	if err != nil {
		t.Fatalf("Erro no Allow: %v", err)
	}
	if allowed {
		t.Fatal("6ª requisição foi permitida, mas deveria ser bloqueada")
	}
}

func TestRateLimiter_Token(t *testing.T) {
	fakeStore := NewFakeStore()
	cfg := &config.Config{
		RateLimitIP:    5,
		RateLimitToken: 10,
		BlockTimeIP:    5 * time.Second,
		BlockTimeToken: 5 * time.Second,
	}
	rl := limiter.NewRateLimiter(fakeStore, cfg)
	ctx := context.Background()
	identifier := "token:abc123"

	for i := 0; i < 10; i++ {
		allowed, err := rl.Allow(ctx, identifier, cfg.RateLimitToken, cfg.BlockTimeToken)
		if err != nil {
			t.Fatalf("Erro no Allow: %v", err)
		}
		if !allowed {
			t.Fatalf("Requisição %d foi bloqueada, mas deveria ser permitida", i+1)
		}
	}

	allowed, err := rl.Allow(ctx, identifier, cfg.RateLimitToken, cfg.BlockTimeToken)
	if err != nil {
		t.Fatalf("Erro no Allow: %v", err)
	}
	if allowed {
		t.Fatal("11ª requisição foi permitida, mas deveria ser bloqueada")
	}
}
