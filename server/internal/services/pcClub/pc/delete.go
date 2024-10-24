package pc

import (
	"context"
	"errors"
	"fmt"
	"server/internal/storage/ssms"
)

func (s *Service) DeletePc(
	ctx context.Context,
	pcId int64,
) error {
	const op = "services.pcClub.pc.DeletePc"

	if err := s.owner.DeletePc(ctx, pcId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
