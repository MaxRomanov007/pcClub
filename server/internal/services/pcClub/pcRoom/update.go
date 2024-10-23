package pcRoom

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/ssms"
)

func (s *Service) UpdatePcRoom(
	ctx context.Context,
	pcRoom models.PcRoom,
) error {
	const op = "services.pcClub.pcRoom.UpdatePcRoom"

	if err := s.owner.UpdatePcRoom(ctx, pcRoom); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to update pc room: %w", op, err)
	}

	return nil
}
