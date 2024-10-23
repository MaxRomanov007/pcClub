package pc

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/ssms"
)

func (s *Service) Pcs(
	ctx context.Context,
	typeId int64,
	isAvailable bool,
) ([]models.PcData, error) {
	const op = "services.pcClub.pc.pcs"

	pcs, err := s.provider.Pcs(ctx, typeId, isAvailable)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return nil, fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return nil, fmt.Errorf("%s: failed to get pcs: %w", op, err)
	}

	return pcs, err
}
