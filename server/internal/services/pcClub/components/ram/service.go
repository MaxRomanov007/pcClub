package ram

import (
	"context"
	"server/internal/models"
)

type provider interface {
	RamTypes(
		ctx context.Context,
	) (ramTypes []models.RAMType, err error)

	Rams(
		ctx context.Context,
		typeID int64,
	) (rams []models.RAM, err error)
}

type owner interface {
	SaveRamType(
		ctx context.Context,
		ramType *models.RAMType,
	) (id int64, err error)

	SaveRam(
		ctx context.Context,
		ram *models.RAM,
	) (id int64, err error)

	DeleteRamType(
		ctx context.Context,
		ramTypeID int64,
	) (err error)

	DeleteRam(
		ctx context.Context,
		ramID int64,
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
	RedisRamTypesKey = "ram_types"
)

func New(provider provider, owner owner, redisProvider redisProvider, redisOwner redisOwner) *Service {
	return &Service{
		provider:      provider,
		owner:         owner,
		redisProvider: redisProvider,
		redisOwner:    redisOwner,
	}
}
