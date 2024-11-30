package ram

import (
	"context"
	"server/internal/lib/errors"
)

func (s *Service) DeleteRamType(
	ctx context.Context,
	typeID int64,
) error {
	const op = "services.pcClub.components.monitor.DeleteRamType"

	if err := s.owner.DeleteRamType(ctx, typeID); err != nil {
		return errors.WithMessage(err, op, "failed to delete monitor type from mssql")
	}

	return nil
}

func (s *Service) DeleteRam(
	ctx context.Context,
	monitorID int64,
) error {
	const op = "services.pcClub.components.monitor.DeleteRam"

	if err := s.owner.DeleteRam(ctx, monitorID); err != nil {
		return errors.WithMessage(err, op, "failed to delete monitor from mssql")
	}

	return nil
}
