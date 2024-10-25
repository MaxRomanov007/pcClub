package dish

import (
	"context"
	"errors"
	"fmt"
	"server/internal/storage/ssms"
)

func (s *Service) DeleteDish(
	ctx context.Context,
	dishId int64,
) error {
	const op = "services.pc.Club.dish.DeleteDish"

	if err := s.owner.DeleteDish(ctx, dishId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete dish: %w", op, err)
	}

	return nil
}
