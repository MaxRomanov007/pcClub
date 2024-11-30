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
		videoCardID int64,
	) (cards []models.VideoCard, err error)
}

type owner interface {
	SaveVideoCardProducer(
		ctx context.Context,
		producer *models.VideoCardProducer,
	) (id int64, err error)

	SaveVideoCard(
		ctx context.Context,
		card *models.VideoCard,
	) (id int64, err error)

	DeleteVideoCardProducer(
		ctx context.Context,
		producerID int64,
	) (err error)

	DeleteVideoCard(
		ctx context.Context,
		cardID int64,
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

const (
	RedisVideoCardProducersKey = "video_card_producers"
)

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
