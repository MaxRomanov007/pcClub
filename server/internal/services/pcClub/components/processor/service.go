package processor

import (
	"golang.org/x/net/context"
	"server/internal/models"
)

type provider interface {
	ProcessorProducers(
		ctx context.Context,
	) (producers []models.ProcessorProducer, err error)

	Processors(
		ctx context.Context,
		producerId int64,
	) (processors []models.Processor, err error)
}

type owner interface {
	SaveProcessorProducer(
		ctx context.Context,
		name string,
	) (err error)

	SaveProcessor(
		ctx context.Context,
		processor models.Processor,
	) (err error)

	DeleteProcessorProducer(
		ctx context.Context,
		producerId int64,
	) (err error)

	DeleteProcessor(
		ctx context.Context,
		processorId int64,
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
