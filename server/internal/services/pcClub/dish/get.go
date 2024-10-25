package dish

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/redis"
	"server/internal/storage/ssms"
)

func (s *Service) Dishes(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]models.DishData, error) {
	const op = "services.pcClub.dish.get.Dishes"

	var dishes []models.DishData
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("dishes:%d-%d", limit, offset),
		&dishes,
	)
	if err == nil {
		return dishes, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get dishes from redis: %w", op, err)
	}

	dishes, err = s.provider.Dishes(ctx, limit, offset)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, HandleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get dishes from ssms: %w", op, err)
	}

	if err := s.redisOwner.Set(
		ctx,
		fmt.Sprintf("dishes:%d-%d", limit, offset),
		dishes,
	); err != nil {
		return nil, fmt.Errorf("%s: failed to write dishes in redis: %w", op, err)
	}

	return dishes, nil
}

func (s *Service) Dish(
	ctx context.Context,
	dishId int64,
) (models.DishData, error) {
	const op = "services.pcClub.dish.get.Dish"

	var dish models.DishData
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("dish:%d", dishId),
		&dish,
	)
	if err == nil {
		return dish, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return models.DishData{}, fmt.Errorf("%s: failed to get dish from redis: %w", op, err)
	}

	dish, err = s.provider.Dish(ctx, dishId)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return models.DishData{}, fmt.Errorf("%s: %w", op, HandleStorageError(ssmsErr))
		}
		return models.DishData{}, fmt.Errorf("%s: failed to get dish from ssms: %w", op, err)
	}

	if err := s.redisOwner.Set(
		ctx,
		fmt.Sprintf("dish:%d", dishId),
		dish,
	); err != nil {
		return models.DishData{}, fmt.Errorf("%s: failed to write dish in redis: %w", op, err)
	}

	return dish, nil
}
