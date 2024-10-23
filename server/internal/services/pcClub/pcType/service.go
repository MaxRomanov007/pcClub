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
		id int64,
	) (pcType models.PcTypeData, err error)

	PcTypes(
		ctx context.Context,
		limit int64,
		offset int64,
	) (pcTypes []models.PcTypeData, err error)
}

type owner interface {
	SavePcType(
		ctx context.Context,
		name string,
		description string,
		processor *models.ProcessorData,
		videoCard *models.VideoCardData,
		monitor *models.MonitorData,
		ram *models.RamData,
	) (err error)

	DeletePcType(
		ctx context.Context,
		typeId int64,
	) (err error)

	UpdatePcType(
		ctx context.Context,
		typeId int64,
		name string,
		description string,
		processor *models.ProcessorData,
		videoCard *models.VideoCardData,
		monitor *models.MonitorData,
		ram *models.RamData,
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
