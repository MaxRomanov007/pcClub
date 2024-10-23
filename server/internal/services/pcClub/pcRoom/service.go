package pcRoom

import (
	"context"
	"server/internal/models"
)

type provider interface {
	PcRoom(
		ctx context.Context,
		pcId int64,
	) (pcRoom models.PcRoom, err error)
}

type owner interface {
	SavePcRoom(
		ctx context.Context,
		pcRoom models.PcRoom,
	) (err error)

	UpdatePcRoom(
		ctx context.Context,
		pcRoom models.PcRoom,
	) (err error)

	DeletePcRoom(
		ctx context.Context,
		roomId int64,
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

func New(redisProvider redisProvider, redisOwner redisOwner, provider provider, owner owner) *Service {
	return &Service{
		redisProvider: redisProvider,
		redisOwner:    redisOwner,
		provider:      provider,
		owner:         owner,
	}
}
