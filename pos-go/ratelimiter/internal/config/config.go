package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	RateLimitIP    int
	RateLimitToken int
	BlockTimeIP    time.Duration
	BlockTimeToken time.Duration
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
}

func LoadConfig() *Config {
	cfg := &Config{}

	if v, exists := os.LookupEnv("RATE_LIMIT_IP"); exists {
		if val, err := strconv.Atoi(v); err == nil {
			cfg.RateLimitIP = val
		} else {
			cfg.RateLimitIP = 5
		}
	} else {
		cfg.RateLimitIP = 5
	}

	if v, exists := os.LookupEnv("RATE_LIMIT_TOKEN"); exists {
		if val, err := strconv.Atoi(v); err == nil {
			cfg.RateLimitToken = val
		} else {
			cfg.RateLimitToken = 10
		}
	} else {
		cfg.RateLimitToken = 10
	}

	if v, exists := os.LookupEnv("BLOCK_TIME_IP"); exists {
		if val, err := strconv.Atoi(v); err == nil {
			cfg.BlockTimeIP = time.Duration(val) * time.Second
		} else {
			cfg.BlockTimeIP = 300 * time.Second
		}
	} else {
		cfg.BlockTimeIP = 300 * time.Second
	}

	if v, exists := os.LookupEnv("BLOCK_TIME_TOKEN"); exists {
		if val, err := strconv.Atoi(v); err == nil {
			cfg.BlockTimeToken = time.Duration(val) * time.Second
		} else {
			cfg.BlockTimeToken = 300 * time.Second
		}
	} else {
		cfg.BlockTimeToken = 300 * time.Second
	}

	if v, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddr = v
	} else {
		cfg.RedisAddr = "localhost:6379"
	}

	if v, exists := os.LookupEnv("REDIS_PASSWORD"); exists {
		cfg.RedisPassword = v
	} else {
		cfg.RedisPassword = ""
	}

	if v, exists := os.LookupEnv("REDIS_DB"); exists {
		if val, err := strconv.Atoi(v); err == nil {
			cfg.RedisDB = val
		} else {
			cfg.RedisDB = 0
		}
	} else {
		cfg.RedisDB = 0
	}

	return cfg
}
