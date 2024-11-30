package ram

import (
	"context"
	"errors"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/redis"
)

func (s *Service) RamTypes(
	ctx context.Context,
) ([]models.RAMType, error) {
	const op = "services.pcClub.components.ram.RamTypes"

	var types []models.RAMType
	err := s.redisProvider.Value(
		ctx,
		RedisRamTypesKey,
		&types,
	)
	if err == nil {
		return types, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, errors2.WithMessage(err, op, "failed to get ram types from redis")
	}

	types, err = s.provider.RamTypes(ctx)
	if err != nil {
		return nil, errors2.WithMessage(components.HandleStorageError(err), op, "failed to get ram types from mssql")
	}

	err = s.redisOwner.Set(
		ctx,
		RedisRamTypesKey,
		types,
	)
	if err != nil {
		return nil, errors2.WithMessage(err, op, "failed to insert ram types into redis")
	}

	return types, nil
}

func (s *Service) Rams(
	ctx context.Context,
	typeID int64,
) ([]models.RAM, error) {
	const op = "services.pcClub.components.ram.RamTypes"

	rams, err := s.provider.Rams(ctx, typeID)
	if err != nil {
		return nil, errors2.WithMessage(
			components.HandleStorageError(err),
			op, "failed to get rams from mssql",
		)
	}

	return rams, nil
}
