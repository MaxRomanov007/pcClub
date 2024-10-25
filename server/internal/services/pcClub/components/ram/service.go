package ram

import (
	"context"
	"server/internal/models"
)

type provider interface {
	RamTypes(
		ctx context.Context,
	) (producers []models.RamType, err error)

	Rams(
		ctx context.Context,
		typeId int64,
	) (rams []models.Ram, err error)
}

type owner interface {
	SaveRamType(
		ctx context.Context,
		name string,
	) (err error)

	SaveRam(
		ctx context.Context,
		ram models.Ram,
	) (err error)

	DeleteRamType(
		ctx context.Context,
		typeId int64,
	) (err error)

	DeleteRam(
		ctx context.Context,
		ramId int64,
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
