package pc

import (
	"context"
	errors2 "server/internal/lib/errors"
)

func (s *Service) DeletePc(
	ctx context.Context,
	pcID int64,
) error {
	const op = "services.pcClub.pc.DeletePc"

	if err := s.owner.DeletePc(ctx, pcID); err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to delete pc from mssql")
	}

	return nil
}
