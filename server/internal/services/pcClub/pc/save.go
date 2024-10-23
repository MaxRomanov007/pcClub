package pc

import (
	"context"
	"errors"
	"fmt"
	"server/internal/storage/ssms"
)

func (s *Service) SavePc(
	ctx context.Context,
	typeId int64,
	roomId int64,
	row int,
	place int,
	description string,
) error {
	const op = "services.pcClub.pc.SavePc"

	err := s.owner.SavePc(
		ctx,
		typeId,
		roomId,
		row,
		place,
		description,
	)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save pc: %w", op, err)
	}

	return nil
}
