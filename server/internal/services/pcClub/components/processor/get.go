package processor

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/redis"
	"server/internal/storage/ssms"
)

func (s *Service) ProcessorProducers(
	ctx context.Context,
) ([]models.ProcessorProducer, error) {
	const op = "services.pcClub.components.processor.ProcessorProducers"

	var producers []models.ProcessorProducer
	err := s.redisProvider.Value(
		ctx,
		"processor_producers",
		&producers,
	)
	if err == nil {
		return producers, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get processor producers from redis: %w", op, err)
	}

	producers, err = s.provider.ProcessorProducers(ctx)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get processor producers: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		"processor_producers",
		producers,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert processor producers into redis: %w", op, err)
	}

	return producers, nil
}

func (s *Service) Processors(
	ctx context.Context,
	producerId int64,
) ([]models.Processor, error) {
	const op = "services.pcClub.components.processor.Processors"

	var processors []models.Processor
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("processors:%d", producerId),
		&processors,
	)
	if err == nil {
		return processors, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get processors from redis: %w", op, err)
	}

	processors, err = s.provider.Processors(ctx, producerId)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get processors: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("processors:%d", producerId),
		processors,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert processors into redis: %w", op, err)
	}

	return processors, nil
}
