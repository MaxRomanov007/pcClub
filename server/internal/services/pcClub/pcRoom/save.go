package pcRoom

import (
	"context"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
)

func (s *Service) SavePcRoom(
	ctx context.Context,
	pcRoom *models.PcRoom,
) (int64, error) {
	const op = "services.pcClub.pcClub.SavePcRoom"

	id, err := s.owner.SavePcRoom(ctx, pcRoom)
	if err != nil {
		return 0, errors2.WithMessage(HandleStorageError(err), op, "failed to save pc room in mssql")
	}

	return id, nil
}
