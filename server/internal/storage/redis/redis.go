package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func (s *Storage) SetStringWithCustomTTL(
	ctx context.Context,
	key string,
	value string,
	ttl time.Duration,
) error {
	const op = "storage.redis.SetStringWithCustomTTL"

	err := s.cl.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: failed to set key: %w", op, err)
	}
	return nil
}

func (s *Storage) Set(
	ctx context.Context,
	key string,
	value interface{},
) error {
	const op = "storage.redis.Set"

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%s: failed to serialize value: %w", op, err)
	}

	err = s.cl.Set(ctx, key, valueJSON, s.cfg.DefaultTTL).Err()
	if err != nil {
		return fmt.Errorf("%s: failed to set key: %w", op, err)
	}

	return nil
}

func (s *Storage) StringValue(
	ctx context.Context,
	key string,
) (string, error) {
	const op = "storage.redis.StringValue"

	value, err := s.cl.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("%s: failed to get value: %w", op, err)
	}

	return value, nil
}

func (s *Storage) Value(
	ctx context.Context,
	key string,
	value interface{},
) error {
	const op = "storage.redis.Value"

	valueJSON, err := s.cl.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: failed to get value: %w", op, err)
	}

	err = json.Unmarshal(valueJSON, value)
	if err != nil {
		return fmt.Errorf("%s: failed to deserialize value: %w", op, err)
	}

	return nil
}
