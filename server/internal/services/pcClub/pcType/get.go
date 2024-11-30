package pcType

import (
	"context"
	"errors"
	"fmt"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
	"server/internal/storage/redis"
)

func (s *Service) PcTypes(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.PcType, error) {
	const op = "services.pcClub.pc.pcTypes"

	var pcTypes []models.PcType
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("pc_types:%d-%d", limit, offset),
		&pcTypes,
	)
	if err == nil {
		return pcTypes, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, errors2.WithMessage(err, op, "failed to get pc types from redis")
	}

	pcTypes, err = s.provider.PcTypes(ctx, limit, offset)
	if err != nil {
		return nil, errors2.WithMessage(HandleStorageError(err), op, "failed to get pc types from mssql")
	}

	err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("pc_types:%d-%d", limit, offset),
		pcTypes,
	)
	if err != nil {
		return nil, errors2.WithMessage(err, op, "failed to save pc types in redis")
	}

	return pcTypes, nil
}

func (s *Service) PcType(
	ctx context.Context,
	typeID int64,
) (models.PcType, error) {
	const op = "services.pcClub.pc.pcType"

	var pcType models.PcType
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("pc_type:%d", typeID),
		&pcType,
	)
	if err == nil {
		return pcType, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return models.PcType{}, errors2.WithMessage(err, op, "failed to get pc type from redis")
	}

	pcType, err = s.provider.PcType(ctx, typeID)
	if err != nil {
		return models.PcType{}, errors2.WithMessage(HandleStorageError(err), op, "failed to get pc type from mssql")
	}

	err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("pc_type:%d", typeID),
		pcType,
	)
	if err != nil {
		return models.PcType{}, errors2.WithMessage(HandleStorageError(err), op, "failed to save pc type in redis")
	}

	return pcType, nil
}
