package processor

import (
	"context"
	"errors"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/redis"
)

func (s *Service) ProcessorProducers(
	ctx context.Context,
) ([]models.ProcessorProducer, error) {
	const op = "services.pcClub.components.processor.ProcessorProducers"

	var producers []models.ProcessorProducer
	err := s.redisProvider.Value(
		ctx,
		RedisProcessorProducersKey,
		&producers,
	)
	if err == nil {
		return producers, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, errors2.WithMessage(err, op, "failed to get processor producers from redis")
	}

	producers, err = s.provider.ProcessorProducers(ctx)
	if err != nil {
		return nil, errors2.WithMessage(components.HandleStorageError(err), op, "failed to get processor producers from mssql")
	}

	err = s.redisOwner.Set(
		ctx,
		RedisProcessorProducersKey,
		producers,
	)
	if err != nil {
		return nil, errors2.WithMessage(err, op, "failed to insert processor producers into redis")
	}

	return producers, nil
}

func (s *Service) Processors(
	ctx context.Context,
	producerID int64,
) ([]models.Processor, error) {
	const op = "services.pcClub.components.processor.ProcessorProducers"

	processors, err := s.provider.Processors(ctx, producerID)
	if err != nil {
		return nil, errors2.WithMessage(
			components.HandleStorageError(err),
			op, "failed to get processors from mssql",
		)
	}

	return processors, nil
}
