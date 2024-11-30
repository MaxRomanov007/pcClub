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
		producerID int64,
	) (processors []models.Processor, err error)
}

type owner interface {
	SaveProcessorProducer(
		ctx context.Context,
		producer *models.ProcessorProducer,
	) (id int64, err error)

	SaveProcessor(
		ctx context.Context,
		processor *models.Processor,
	) (id int64, err error)

	DeleteProcessorProducer(
		ctx context.Context,
		producerID int64,
	) (err error)

	DeleteProcessor(
		ctx context.Context,
		processorID int64,
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

const (
	RedisProcessorProducersKey = "processor_producers"
)

func New(provider provider, owner owner, redisProvider redisProvider, redisOwner redisOwner) *Service {
	return &Service{
		provider:      provider,
		owner:         owner,
		redisProvider: redisProvider,
		redisOwner:    redisOwner,
	}
}
