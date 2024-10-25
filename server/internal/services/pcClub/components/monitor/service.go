package monitor

import (
	"context"
	"server/internal/models"
)

type provider interface {
	MonitorProducers(
		ctx context.Context,
	) (producers []models.MonitorProducer, err error)

	Monitors(
		ctx context.Context,
		producerId int64,
	) (monitors []models.Monitor, err error)
}

type owner interface {
	SaveMonitorProducer(
		ctx context.Context,
		name string,
	) (err error)

	SaveMonitor(
		ctx context.Context,
		monitor models.Monitor,
	) (err error)

	DeleteMonitorProducer(
		ctx context.Context,
		producerId int64,
	) (err error)

	DeleteMonitor(
		ctx context.Context,
		monitorId int64,
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
