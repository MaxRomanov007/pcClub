package ram

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/redis"
	"server/internal/storage/ssms"
)

func (s *Service) RamTypes(
	ctx context.Context,
) ([]models.RamType, error) {
	const op = "services.pcClub.components.ram.RamTypes"

	var types []models.RamType
	err := s.redisProvider.Value(
		ctx,
		"ram_types",
		&types,
	)
	if err == nil {
		return types, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get ram types from redis: %w", op, err)
	}

	types, err = s.provider.RamTypes(ctx)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get ram types: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		"ram_types",
		types,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert ram types into redis: %w", op, err)
	}

	return types, nil
}

func (s *Service) Rams(
	ctx context.Context,
	typeId int64,
) ([]models.Ram, error) {
	const op = "services.pcClub.components.ram.Rams"

	var rams []models.Ram
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("rams:%d", typeId),
		&rams,
	)
	if err == nil {
		return rams, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get rams from redis: %w", op, err)
	}

	rams, err = s.provider.Rams(ctx, typeId)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get rams: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("rams:%d", typeId),
		rams,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert rams into redis: %w", op, err)
	}

	return rams, nil
}
