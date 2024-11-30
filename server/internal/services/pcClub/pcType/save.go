package pcType

import (
	"context"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
)

func (s *Service) SavePcType(
	ctx context.Context,
	pcType *models.PcType,
) (int64, error) {
	const op = "services.pcClub.pc.SavePcType"

	id, err := s.owner.SavePcType(ctx, pcType)
	if err != nil {
		return 0, errors2.WithMessage(HandleStorageError(err), op, "failed to save pc type in mssql")
	}

	return id, nil
}
