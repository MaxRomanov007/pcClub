package pc

import (
	"context"
	"server/internal/models"
)

type provider interface {
	Pcs(
		ctx context.Context,
		typeID int64,
		isAvailable bool,
	) (pcs []models.Pc, err error)
}

type owner interface {
	SavePc(
		ctx context.Context,
		pc *models.Pc,
	) (id int64, err error)

	UpdatePc(
		ctx context.Context,
		pcID int64,
		pc *models.Pc,
	) (err error)

	DeletePc(
		ctx context.Context,
		pcID int64,
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
