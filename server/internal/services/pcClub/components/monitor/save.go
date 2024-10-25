package monitor

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/ssms"
)

func (s *Service) SaveMonitorProducer(
	ctx context.Context,
	name string,
) error {
	const op = "services.pcClub.components.monitor.SaveMonitorProducer"

	if err := s.owner.SaveMonitorProducer(ctx, name); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save monitor producer: %w", op, err)
	}

	return nil
}

func (s *Service) SaveMonitor(
	ctx context.Context,
	monitor models.Monitor,
) error {
	const op = "services.pcClub.components.monitor.SaveMonitor"

	if err := s.owner.SaveMonitor(ctx, monitor); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save monitor: %w", op, err)
	}

	return nil
}
