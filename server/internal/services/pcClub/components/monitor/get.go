package monitor

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/redis"
	"server/internal/storage/ssms"
)

func (s *Service) MonitorProducers(
	ctx context.Context,
) ([]models.MonitorProducer, error) {
	const op = "services.pcClub.components.monitor.MonitorProducers"

	var producers []models.MonitorProducer
	err := s.redisProvider.Value(
		ctx,
		"monitor_producers",
		&producers,
	)
	if err == nil {
		return producers, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get monitor producers from redis: %w", op, err)
	}

	producers, err = s.provider.MonitorProducers(ctx)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get monitor producers: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		"monitor_producers",
		producers,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert monitor producers into redis: %w", op, err)
	}

	return producers, nil
}

func (s *Service) Monitors(
	ctx context.Context,
	producerId int64,
) ([]models.Monitor, error) {
	const op = "services.pcClub.components.monitor.Monitors"

	var monitors []models.Monitor
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("monitors:%d", producerId),
		&monitors,
	)
	if err == nil {
		return monitors, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get monitors from redis: %w", op, err)
	}

	monitors, err = s.provider.Monitors(ctx, producerId)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get monitors: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("monitors:%d", producerId),
		monitors,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert monitors into redis: %w", op, err)
	}

	return monitors, nil
}
