package dish

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/ssms"
)

func (s *Service) SaveDish(
	ctx context.Context,
	dish models.DishData,
) error {
	const op = "services.pc.Club.dish.SaveDish"

	if err := s.owner.SaveDish(ctx, dish); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save dish: %w", op, err)
	}

	return nil
}
