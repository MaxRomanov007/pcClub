package pcType

import (
	"context"
	"server/internal/models"
)

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

type provider interface {
	PcType(
		ctx context.Context,
		typeID int64,
	) (pcType models.PcType, err error)

	PcTypes(
		ctx context.Context,
		limit int,
		offset int,
	) (pcTypes []models.PcType, err error)
}

type owner interface {
	SavePcType(
		ctx context.Context,
		pcType *models.PcType,
	) (id int64, err error)

	UpdatePcType(
		ctx context.Context,
		typeID int64,
		pcType *models.PcType,
	) (err error)

	DeletePcType(
		ctx context.Context,
		typeId int64,
	) (err error)
}

type Service struct {
	redisProvider redisProvider
	redisOwner    redisOwner
	provider      provider
	owner         owner
}

func New(provider provider, owner owner, redisProvider redisProvider, redisOwner redisOwner) *Service {
	return &Service{
		redisProvider: redisProvider,
		redisOwner:    redisOwner,
		provider:      provider,
		owner:         owner,
	}
}
