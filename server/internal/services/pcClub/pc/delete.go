package pc

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
		if errors.Is(err, ssms.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return fmt.Errorf("%s: failed to delete pc type: %w", op, err)
	}

	return nil
}
