package dish

import (
	"context"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
)

func (s *Service) SaveDish(
	ctx context.Context,
	dish *models.Dish,
) (int64, error) {
	const op = "services.pc.Club.dish.SaveDish"

	id, err := s.owner.SaveDish(ctx, dish)
	if err != nil {
		return 0, errors2.WithMessage(HandleStorageError(err), op, "failed to save dish in mssql")
	}

	return id, nil
}
