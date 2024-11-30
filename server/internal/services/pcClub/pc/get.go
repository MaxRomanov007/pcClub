package pc

import (
	"context"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
)

func (s *Service) Pcs(
	ctx context.Context,
	typeId int64,
	isAvailable bool,
) ([]models.Pc, error) {
	const op = "services.pcClub.pc.pcs"

	pcs, err := s.provider.Pcs(ctx, typeId, isAvailable)
	if err != nil {
		return nil, errors2.WithMessage(HandleStorageError(err), op, "failed to get pcs from mssql")
	}

	return pcs, nil
}
