package monitor

import (
	"context"
	"errors"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/redis"
)

func (s *Service) MonitorProducers(
	ctx context.Context,
) ([]models.MonitorProducer, error) {
	const op = "services.pcClub.components.monitor.MonitorProducers"

	var producers []models.MonitorProducer
	err := s.redisProvider.Value(
		ctx,
		RedisMonitorProducersKey,
		&producers,
	)
	if err == nil {
		return producers, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, errors2.WithMessage(err, op, "failed to get monitor producers from redis")
	}

	producers, err = s.provider.MonitorProducers(ctx)
	if err != nil {
		return nil, errors2.WithMessage(components.HandleStorageError(err), op, "failed to get monitor producers from mssql")
	}

	err = s.redisOwner.Set(
		ctx,
		RedisMonitorProducersKey,
		producers,
	)
	if err != nil {
		return nil, errors2.WithMessage(err, op, "failed to insert monitor producers into redis")
	}

	return producers, nil
}

func (s *Service) Monitors(
	ctx context.Context,
	producerID int64,
) ([]models.Monitor, error) {
	const op = "services.pcClub.components.monitor.MonitorProducers"

	monitors, err := s.provider.Monitors(ctx, producerID)
	if err != nil {
		return nil, errors2.WithMessage(
			components.HandleStorageError(err),
			op, "failed to get monitors from mssql",
		)
	}

	return monitors, nil
}
