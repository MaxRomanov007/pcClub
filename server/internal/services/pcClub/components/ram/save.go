package ram

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/ssms"
)

func (s *Service) SaveRamType(
	ctx context.Context,
	name string,
) error {
	const op = "services.pcClub.components.ram.SaveRamType"

	if err := s.owner.SaveRamType(ctx, name); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save ram type: %w", op, err)
	}

	return nil
}

func (s *Service) SaveRam(
	ctx context.Context,
	ram models.Ram,
) error {
	const op = "services.pcClub.components.ram.SaveRam"

	if err := s.owner.SaveRam(ctx, ram); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save ram: %w", op, err)
	}

	return nil
}
