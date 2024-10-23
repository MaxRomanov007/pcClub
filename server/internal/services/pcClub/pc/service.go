package pc

import (
	"context"
	"server/internal/models"
)

type provider interface {
	Pcs(
		ctx context.Context,
		typeId int64,
		isAvailable bool,
	) (pcs []models.PcData, err error)
}

type owner interface {
	SavePc(
		ctx context.Context,
		typeId int64,
		roomId int64,
		row int,
		place int,
		description string,
	) (err error)

	DeletePc(
		ctx context.Context,
		pcId int64,
	) (err error)

	UpdatePc(
		ctx context.Context,
		pcId int64,
		typeId int64,
		roomId int64,
		statusId int64,
		row int,
		place int,
		description string,
	) (err error)
}

type Service struct {
	provider provider
	owner    owner
}

func New(provider provider, owner owner) *Service {
	return &Service{
		provider: provider,
		owner:    owner,
	}
}
