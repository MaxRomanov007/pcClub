package monitor

import (
	"context"
	"errors"
	"fmt"
	"server/internal/services/pcClub/components"
	"server/internal/storage/ssms"
)

func (s *Service) DeleteMonitorProducer(
	ctx context.Context,
	producerId int64,
) error {
	const op = "services.pcClub.components.monitor.DeleteMonitorProducer"

	if err := s.owner.DeleteMonitorProducer(ctx, producerId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete monitor producer: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteMonitor(
	ctx context.Context,
	monitorId int64,
) error {
	const op = "services.pcClub.components.monitor.DeleteMonitor"

	if err := s.owner.DeleteMonitor(ctx, monitorId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete monitor: %w", op, err)
	}

	return nil
}
