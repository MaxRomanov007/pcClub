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
		producerID int64,
	) (monitors []models.Monitor, err error)
}

type owner interface {
	SaveMonitorProducer(
		ctx context.Context,
		producer *models.MonitorProducer,
	) (id int64, err error)

	SaveMonitor(
		ctx context.Context,
		monitor *models.Monitor,
	) (id int64, err error)

	DeleteMonitorProducer(
		ctx context.Context,
		producerID int64,
	) (err error)

	DeleteMonitor(
		ctx context.Context,
		monitorID int64,
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
	RedisMonitorProducersKey = "monitor_producers"
)

func New(provider provider, owner owner, redisProvider redisProvider, redisOwner redisOwner) *Service {
	return &Service{
		provider:      provider,
		owner:         owner,
		redisProvider: redisProvider,
		redisOwner:    redisOwner,
	}
}
