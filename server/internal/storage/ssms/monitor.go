package ssms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

func (s *Storage) MonitorProducers(
	ctx context.Context,
) ([]models.MonitorProducer, error) {
	const op = "storage.ssms.monitor.MonitorProducers"

	query, args, err := squirrel.Select("*").From("monitor_producers").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var producers []models.MonitorProducer
	if err := s.db.SelectContext(ctx, &producers, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to select monitor producers: %w", op, handleError(err))
	}

	return producers, nil
}

func (s *Storage) Monitors(
	ctx context.Context,
	producerId int64,
) ([]models.Monitor, error) {
	const op = "storage.ssms.monitor.Monitors"

	query, args, err := squirrel.
		Select("*").
		From("monitors").
		Where(squirrel.Eq{
			"monitor_producer_id": producerId,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var monitors []models.Monitor
	if err := s.db.SelectContext(ctx, &monitors, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to select monitors: %w", op, handleError(err))
	}

	return monitors, nil
}

func (s *Storage) SaveMonitorProducer(
	ctx context.Context,
	name string,
) error {
	const op = "storage.ssms.monitor.SaveMonitorProducer"

	query, args, err := squirrel.Insert("monitor_producers").Values(name).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to save monitor producer: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) SaveMonitor(
	ctx context.Context,
	monitor models.Monitor,
) error {
	const op = "storage.ssms.monitor.SaveMonitor"

	query, args, err := squirrel.
		Insert("monitors").
		Columns(
			"monitor_producer_id",
			"model",
		).
		Values(
			monitor.MonitorProducerID,
			monitor.Model,
		).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to save monitor: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeleteMonitorProducer(
	ctx context.Context,
	producerId int64,
) error {
	const op = "storage.ssms.monitor.DeleteMonitorProducer"

	query, args, err := squirrel.
		Delete("monitor_producers").
		Where(squirrel.Eq{
			"monitor_producer_id": producerId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to delete monitor producer: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeleteMonitor(
	ctx context.Context,
	monitorId int64,
) error {
	const op = "storage.ssms.monitor.DeleteMonitor"

	query, args, err := squirrel.
		Delete("monitors").
		Where(squirrel.Eq{
			"monitor_id": monitorId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to delete monitor: %w", op, handleError(err))
	}

	return nil
}
