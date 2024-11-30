package dish

import (
	"context"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
)

func (s *Service) UpdateDish(
	ctx context.Context,
	dishID int64,
	dish *models.Dish,
) error {
	const op = "services.pc.Club.dish.UpdateDish"

	if err := s.owner.UpdateDish(ctx, dishID, dish); err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to update dish in mssql")
	}

	return nil
}
