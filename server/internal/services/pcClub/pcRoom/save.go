package pcRoom

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/ssms"
)

func (s *Service) SavePcRoom(
	ctx context.Context,
	pcRoom models.PcRoom,
) error {
	const op = "services.pcClub.pcClub.SavePcRoom"

	if err := s.owner.SavePcRoom(ctx, pcRoom); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save pc room: %w", op, err)
	}

	return nil
}
