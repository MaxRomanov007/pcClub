package videoCard

import (
	"context"
	"errors"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/redis"
)

func (s *Service) VideoCardProducers(
	ctx context.Context,
) ([]models.VideoCardProducer, error) {
	const op = "services.pcClub.components.videoCard.VideoCardProducers"

	var producers []models.VideoCardProducer
	err := s.redisProvider.Value(
		ctx,
		RedisVideoCardProducersKey,
		&producers,
	)
	if err == nil {
		return producers, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, errors2.WithMessage(err, op, "failed to get video card producers from redis")
	}

	producers, err = s.provider.VideoCardProducers(ctx)
	if err != nil {
		return nil, errors2.WithMessage(components.HandleStorageError(err), op, "failed to get video card producers from mssql")
	}

	err = s.redisOwner.Set(
		ctx,
		RedisVideoCardProducersKey,
		producers,
	)
	if err != nil {
		return nil, errors2.WithMessage(err, op, "failed to insert video card producers into redis")
	}

	return producers, nil
}

func (s *Service) VideoCards(
	ctx context.Context,
	producerID int64,
) ([]models.VideoCard, error) {
	const op = "services.pcClub.components.videoCard.VideoCardProducers"

	cards, err := s.provider.VideoCards(ctx, producerID)
	if err != nil {
		return nil, errors2.WithMessage(
			components.HandleStorageError(err),
			op, "failed to get video cards from mssql",
		)
	}

	return cards, nil
}
