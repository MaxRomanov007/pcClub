package pc

import (
	"context"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
)

func (s *Service) UpdatePc(
	ctx context.Context,
	pcID int64,
	pc *models.Pc,
) error {
	const op = "services.pcClub.pc.UpdatePc"

	if err := s.owner.UpdatePc(ctx, pcID, pc); err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to update pc in mssql")
	}

	return nil
}
