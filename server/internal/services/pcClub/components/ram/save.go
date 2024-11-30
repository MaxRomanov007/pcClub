package ram

import (
	"context"
	"server/internal/lib/errors"
	"server/internal/models"
	"server/internal/services/pcClub/components"
)

func (s *Service) SaveRamType(
	ctx context.Context,
	ramType *models.RAMType,
) (int64, error) {
	const op = "services.pcClub.components.ram.SaveRamType"

	id, err := s.owner.SaveRamType(ctx, ramType)
	if err != nil {
		return 0, errors.WithMessage(components.HandleStorageError(err), op, "failed to save ram type in mssql")
	}

	return id, nil
}

func (s *Service) SaveRam(
	ctx context.Context,
	ram *models.RAM,
) (int64, error) {
	const op = "services.pcClub.components.ram.SaveRam"

	id, err := s.owner.SaveRam(ctx, ram)
	if err != nil {
		return 0, errors.WithMessage(components.HandleStorageError(err), op, "failed to save ram in mssql")
	}

	return id, nil
}
