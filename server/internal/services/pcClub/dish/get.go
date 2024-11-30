package dish

import (
	"context"
	"errors"
	"fmt"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
	"server/internal/storage/redis"
)

func (s *Service) Dishes(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Dish, error) {
	const op = "services.pcClub.dish.get.Dishes"

	var dishes []models.Dish
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("dishes:%d-%d", limit, offset),
		&dishes,
	)
	if err == nil {
		return dishes, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, errors2.WithMessage(err, op, "failed to get dishes from redis")
	}

	dishes, err = s.provider.Dishes(ctx, limit, offset)
	if err != nil {
		return nil, errors2.WithMessage(HandleStorageError(err), op, "failed to get dishes from mssql")
	}

	if err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("dishes:%d-%d", limit, offset),
		dishes,
	); err != nil {
		return nil, errors2.WithMessage(err, op, "failed to write dishes in redis")
	}

	return dishes, nil
}

func (s *Service) Dish(
	ctx context.Context,
	dishID int64,
) (models.Dish, error) {
	const op = "services.pcClub.dish.get.Dish"

	var dish models.Dish
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("dish:%d", dishID),
		&dish,
	)
	if err == nil {
		return dish, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return models.Dish{}, errors2.WithMessage(err, op, "failed to get dish from redis")
	}

	dish, err = s.provider.Dish(ctx, dishID)
	if err != nil {
		return models.Dish{}, errors2.WithMessage(HandleStorageError(err), op, "failed to get dishes from mssql")
	}

	if err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("dish:%d", dishID),
		dish,
	); err != nil {
		return models.Dish{}, errors2.WithMessage(err, op, "failed to write dish in redis")
	}

	return dish, nil
}
