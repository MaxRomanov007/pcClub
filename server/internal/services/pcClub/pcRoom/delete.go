package pcRoom

import (
	"context"
	errors2 "server/internal/lib/errors"
)

func (s *Service) DeletePcRoom(
	ctx context.Context,
	roomID int64,
) error {
	const op = "services.pcClub.pcRoom.DeletePcRoom"

	if err := s.owner.DeletePcRoom(ctx, roomID); err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to delete pc room from mssql")
	}

	return nil
}
