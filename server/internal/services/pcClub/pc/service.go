package pc

import (
	"context"
	"errors"
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
	Pcs(
		ctx context.Context,
		typeId int64,
		isAvailable bool,
	) (pcs []models.PcData, err error)

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
	SavePc(
		ctx context.Context,
		typeId int64,
		roomId int64,
		row int,
		place int,
	) (err error)

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

	//UpdatePc(
	//	ctx context.Context,
	//	id int64,
	//	pc models.Pc,
	//) (err error)
}

type Service struct {
	redisProvider redisProvider
	redisOwner    redisOwner
	provider      provider
	owner         owner
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("pc type already exists")
)

func New(redisProvider redisProvider, redisOwner redisOwner, provider provider, owner owner) *Service {
	return &Service{
		redisProvider: redisProvider,
		redisOwner:    redisOwner,
		provider:      provider,
		owner:         owner,
	}
}
