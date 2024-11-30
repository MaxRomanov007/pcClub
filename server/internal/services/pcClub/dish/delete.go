package dish

import (
	"context"
	errors2 "server/internal/lib/errors"
)

func (s *Service) DeleteDish(
	ctx context.Context,
	dishId int64,
) error {
	const op = "services.pc.Club.dish.DeleteDish"

	if err := s.owner.DeleteDish(ctx, dishId); err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to delete dish from mssql")
	}

	return nil
}
