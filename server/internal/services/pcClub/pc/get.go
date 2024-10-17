package pc

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/redis"
	"server/internal/storage/ssms"
)

func (s *Service) Pcs(
	ctx context.Context,
	typeId int64,
	isAvailable bool,
) ([]models.PcData, error) {
	const op = "services.pcClub.pc.pcs"

	pcs, err := s.provider.Pcs(ctx, typeId, isAvailable)
	if errors.Is(err, ssms.ErrNotFound) {
		return nil, fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get pcs: %w", op, err)
	}

	return pcs, err
}

func (s *Service) PcTypes(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]models.PcTypeData, error) {
	const op = "services.pcClub.pc.pcTypes"

	var pcTypes []models.PcTypeData
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("pc_types:%d-%d", limit, offset),
		&pcTypes,
	)
	if err == nil {
		return pcTypes, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get pc types from redis: %w", op, err)
	}

	pcTypes, err = s.provider.PcTypes(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get pc types from sql: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("pc_types:%d-%d", limit, offset),
		pcTypes,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert pc types into redis: %w", op, err)
	}

	return pcTypes, nil
}

func (s *Service) PcType(
	ctx context.Context,
	typeId int64,
) (models.PcTypeData, error) {
	const op = "services.pcClub.pc.pcType"

	var pcType models.PcTypeData
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("pc_type:%d", typeId),
		&pcType,
	)
	if err == nil {
		return pcType, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return models.PcTypeData{}, fmt.Errorf("%s: failed to get pc type from redis: %w", op, err)
	}

	pcType, err = s.provider.PcType(ctx, typeId)
	if errors.Is(err, ssms.ErrNotFound) {
		return models.PcTypeData{}, fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	if err != nil {
		return models.PcTypeData{}, fmt.Errorf("%s: failed to get pc type from sql: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("pc_type:%d", typeId),
		pcType,
	)
	if err != nil {
		return models.PcTypeData{}, fmt.Errorf("%s: failed to insert pc type into redis: %w", op, err)
	}

	return pcType, nil
}
