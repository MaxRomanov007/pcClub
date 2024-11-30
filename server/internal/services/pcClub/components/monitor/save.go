package monitor

import (
	"context"
	"server/internal/lib/errors"
	"server/internal/models"
	"server/internal/services/pcClub/components"
)

func (s *Service) SaveMonitorProducer(
	ctx context.Context,
	producer *models.MonitorProducer,
) (int64, error) {
	const op = "services.pcClub.components.monitor.SaveMonitorProducer"

	id, err := s.owner.SaveMonitorProducer(ctx, producer)
	if err != nil {
		return 0, errors.WithMessage(components.HandleStorageError(err), op, "failed to save monitor producer in mssql")
	}

	return id, nil
}

func (s *Service) SaveMonitor(
	ctx context.Context,
	monitor *models.Monitor,
) (int64, error) {
	const op = "services.pcClub.components.monitor.SaveMonitor"

	id, err := s.owner.SaveMonitor(ctx, monitor)
	if err != nil {
		return 0, errors.WithMessage(components.HandleStorageError(err), op, "failed to save monitor in mssql")
	}

	return id, nil
}
