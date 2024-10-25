package videoCard

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/redis"
	"server/internal/storage/ssms"
)

func (s *Service) VideoCardProducers(
	ctx context.Context,
) ([]models.VideoCardProducer, error) {
	const op = "services.pcClub.components.videoCard.VideoCardProducers"

	var producers []models.VideoCardProducer
	err := s.redisProvider.Value(
		ctx,
		"videoCard_producers",
		&producers,
	)
	if err == nil {
		return producers, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get video card producers from redis: %w", op, err)
	}

	producers, err = s.provider.VideoCardProducers(ctx)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get video card producers: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		"videoCard_producers",
		producers,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert video card producers into redis: %w", op, err)
	}

	return producers, nil
}

func (s *Service) VideoCards(
	ctx context.Context,
	producerId int64,
) ([]models.VideoCard, error) {
	const op = "services.pcClub.components.videoCard.VideoCards"

	var videoCards []models.VideoCard
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("videoCards:%d", producerId),
		&videoCards,
	)
	if err == nil {
		return videoCards, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, fmt.Errorf("%s: failed to get video cards from redis: %w", op, err)
	}

	videoCards, err = s.provider.VideoCards(ctx, producerId)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get video cards: %w", op, err)
	}

	err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("videoCards:%d", producerId),
		videoCards,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert video cards into redis: %w", op, err)
	}

	return videoCards, nil
}
