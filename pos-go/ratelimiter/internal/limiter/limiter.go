package limiter

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"rate-limiter/internal/config"
	"rate-limiter/internal/store"
)

const (
	headerAPIKey = "API_KEY"
	errorMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

type RateLimiter struct {
	store  store.RateLimiterStore
	config *config.Config
}

func NewRateLimiter(store store.RateLimiterStore, cfg *config.Config) *RateLimiter {
	return &RateLimiter{
		store:  store,
		config: cfg,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := r.Header.Get(headerAPIKey)
		var identifier string
		var limit int
		var blockTime time.Duration

		if token != "" {
			identifier = "token:" + strings.TrimSpace(token)
			limit = rl.config.RateLimitToken
			blockTime = rl.config.BlockTimeToken
		} else {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			identifier = "ip:" + ip
			limit = rl.config.RateLimitIP
			blockTime = rl.config.BlockTimeIP
		}

		allowed, err := rl.Allow(ctx, identifier, limit, blockTime)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(w, errorMessage, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) Allow(ctx context.Context, identifier string, limit int, blockTime time.Duration) (bool, error) {
	blocked, err := rl.store.IsBlocked(ctx, identifier)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	window := time.Now().Unix()
	key := fmt.Sprintf("rl:%s:%d", identifier, window)

	count, err := rl.store.Increment(ctx, key, time.Second)
	if err != nil {
		return false, err
	}

	if count > int64(limit) {
		err = rl.store.Block(ctx, identifier, blockTime)
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}
