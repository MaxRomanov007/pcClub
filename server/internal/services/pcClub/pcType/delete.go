package pcType

import (
	"context"
	errors2 "server/internal/lib/errors"
)

func (s *Service) DeletePcType(
	ctx context.Context,
	typeID int64,
) error {
	const op = "services.pcClub.pc.DeletePcType"

	if err := s.owner.DeletePcType(ctx, typeID); err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to delete pc type from mssql")
	}

	return nil
}
