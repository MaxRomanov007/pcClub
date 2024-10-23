package pcRoom

import (
	"context"
	"errors"
	"fmt"
	"server/internal/storage/ssms"
)

func (s *Service) DeletePcRoom(
	ctx context.Context,
	roomId int64,
) error {
	const op = "services.pcClub.pcRoom.DeletePcRoom"

	if err := s.owner.DeletePcRoom(ctx, roomId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", handleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
