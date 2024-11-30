package pcRoom

import (
	"context"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
)

func (s *Service) UpdatePcRoom(
	ctx context.Context,
	roomID int64,
	pcRoom *models.PcRoom,
) error {
	const op = "services.pcClub.pcRoom.UpdatePcRoom"

	if err := s.owner.UpdatePcRoom(ctx, roomID, pcRoom); err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to update pc room in mssql")
	}

	return nil
}
