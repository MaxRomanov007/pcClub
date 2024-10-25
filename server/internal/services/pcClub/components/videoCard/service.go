package videoCard

import (
	"context"
	"server/internal/models"
)

type provider interface {
	VideoCardProducers(
		ctx context.Context,
	) (producers []models.VideoCardProducer, err error)

	VideoCards(
		ctx context.Context,
		producerId int64,
	) (videoCards []models.VideoCard, err error)
}

type owner interface {
	SaveVideoCardProducer(
		ctx context.Context,
		name string,
	) (err error)

	SaveVideoCard(
		ctx context.Context,
		videoCard models.VideoCard,
	) (err error)

	DeleteVideoCardProducer(
		ctx context.Context,
		producerId int64,
	) (err error)

	DeleteVideoCard(
		ctx context.Context,
		videoCardId int64,
	) (err error)
}

type redisProvider interface {
	Value(
		ctx context.Context,
		key string,
		value interface{},
	) (err error)
}

type redisOwner interface {
	Set(
		ctx context.Context,
		key string,
		value interface{},
	) (err error)
}

type Service struct {
	provider      provider
	owner         owner
	redisProvider redisProvider
	redisOwner    redisOwner
}

func New(provider provider, owner owner, redisProvider redisProvider, redisOwner redisOwner) *Service {
	return &Service{
		provider:      provider,
		owner:         owner,
		redisProvider: redisProvider,
		redisOwner:    redisOwner,
	}
}
