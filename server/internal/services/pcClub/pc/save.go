package pc

import (
	"context"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
)

func (s *Service) SavePc(
	ctx context.Context,
	pc *models.Pc,
) (int64, error) {
	const op = "services.pcClub.pc.SavePc"

	id, err := s.owner.SavePc(ctx, pc)
	if err != nil {
		return 0, errors2.WithMessage(HandleStorageError(err), op, "failed to save pc in mssql")
	}

	return id, nil
}
