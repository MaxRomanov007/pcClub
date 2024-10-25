package dish

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/ssms"
)

func (s *Service) UpdateDish(
	ctx context.Context,
	dish models.DishData,
) error {
	const op = "services.pc.Club.dish.UpdateDish"

	if err := s.owner.UpdateDish(ctx, dish); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to update dish: %w", op, err)
	}

	return nil
}
