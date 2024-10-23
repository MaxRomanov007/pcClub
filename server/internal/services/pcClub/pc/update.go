package pc

import (
	"context"
	"errors"
	"fmt"
	"server/internal/storage/ssms"
)

func (s *Service) UpdatePc(
	ctx context.Context,
	pcId int64,
	typeId int64,
	roomId int64,
	statusId int64,
	row int,
	place int,
	description string,
) error {
	const op = "services.pcClub.pc.UpdatePc"

	if err := s.owner.UpdatePc(
		ctx,
		pcId,
		typeId,
		roomId,
		statusId,
		row,
		place,
		description,
	); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to update pc: %w", op, err)
	}

	return nil
}
