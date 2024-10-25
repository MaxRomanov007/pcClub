package ram

import (
	"context"
	"errors"
	"fmt"
	"server/internal/services/pcClub/components"
	"server/internal/storage/ssms"
)

func (s *Service) DeleteRamType(
	ctx context.Context,
	typeId int64,
) error {
	const op = "services.pcClub.components.ram.DeleteRamType"

	if err := s.owner.DeleteRamType(ctx, typeId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete ram type: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteRam(
	ctx context.Context,
	ramId int64,
) error {
	const op = "services.pcClub.components.ram.DeleteRam"

	if err := s.owner.DeleteRam(ctx, ramId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete ram: %w", op, err)
	}

	return nil
}
