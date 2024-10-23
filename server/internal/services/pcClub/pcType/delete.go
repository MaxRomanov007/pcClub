package pcType

import (
	"context"
	"errors"
	"fmt"
	"server/internal/storage/ssms"
)

func (s *Service) DeletePcType(
	ctx context.Context,
	typeId int64,
) error {
	const op = "services.pcClub.pc.DeletePcType"

	if err := s.owner.DeletePcType(ctx, typeId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete pc type: %w", op, err)
	}

	return nil
}
