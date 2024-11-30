package monitor

import (
	"context"
	"server/internal/lib/errors"
)

func (s *Service) DeleteMonitorProducer(
	ctx context.Context,
	producerID int64,
) error {
	const op = "services.pcClub.components.monitor.DeleteMonitorProducer"

	if err := s.owner.DeleteMonitorProducer(ctx, producerID); err != nil {
		return errors.WithMessage(err, op, "failed to delete monitor producer from mssql")
	}

	return nil
}

func (s *Service) DeleteMonitor(
	ctx context.Context,
	monitorID int64,
) error {
	const op = "services.pcClub.components.monitor.DeleteMonitor"

	if err := s.owner.DeleteMonitor(ctx, monitorID); err != nil {
		return errors.WithMessage(err, op, "failed to delete monitor from mssql")
	}

	return nil
}
