package pcType

import (
	"context"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
)

func (s *Service) UpdatePcType(
	ctx context.Context,
	typeID int64,
	pcType *models.PcType,
) error {
	const op = "services.pcClub.pc.UpdatePcType"

	if err := s.owner.UpdatePcType(ctx, typeID, pcType); err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to update pc type in mssql")
	}

	return nil
}
