package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"server/internal/config"
)

var (
	ErrNotFound = errors.New("key not found")
)

type Storage struct {
	cfg *config.RedisConfig
	cl  *redis.Client
}

func New(ctx context.Context, cfg *config.RedisConfig) (*Storage, error) {
	const op = "storage.redis.New"

	cl := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	_, err := cl.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect redis: %w", op, err)
	}

	return &Storage{
		cfg: cfg,
		cl:  cl,
	}, nil
}
