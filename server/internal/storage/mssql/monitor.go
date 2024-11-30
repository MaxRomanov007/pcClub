package mssql

import (
	"context"
	"server/internal/lib/api/database/gorm"
	"server/internal/lib/errors"
	"server/internal/models"
)

func (s *Storage) MonitorProducers(
	ctx context.Context,
) ([]models.MonitorProducer, error) {
	const op = "storage.mssql.monitor.MonitorProducers"

	var producers []models.MonitorProducer
	if res := s.db.WithContext(ctx).Find(&producers); gorm.IsFailResult(res) {
		return nil, errors.WithMessage(errorByResult(res), op, "failed to get monitor producers")
	}

	return producers, nil
}

func (s *Storage) Monitors(
	ctx context.Context,
	producerID int64,
) ([]models.Monitor, error) {
	const op = "storage.mssql.monitor.MonitorProducer"

	var monitors []models.Monitor
	if res := s.db.WithContext(ctx).
		Where("monitor_producer_id = ?", producerID).
		Find(&monitors); gorm.IsFailResult(res) {

		return nil, errors.WithMessage(
			errorByResult(res),
			op, "failed to get monitors",
		)
	}

	return monitors, nil
}

func (s *Storage) SaveMonitorProducer(
	ctx context.Context,
	producer *models.MonitorProducer,
) (int64, error) {
	const op = "storage.mssql.monitor.SaveMonitorProducer"

	if res := s.db.WithContext(ctx).Save(producer); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save monitor producer")
	}

	return producer.MonitorProducerID, nil
}

func (s *Storage) SaveMonitor(
	ctx context.Context,
	monitor *models.Monitor,
) (int64, error) {
	const op = "storage.mssql.monitor.SaveMonitor"

	if res := s.db.WithContext(ctx).Save(monitor); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save monitor")
	}

	return monitor.MonitorID, nil
}

func (s *Storage) DeleteMonitorProducer(
	ctx context.Context,
	producerID int64,
) error {
	const op = "storage.mssql.monitor.DeleteMonitorProducer"

	if res := s.db.WithContext(ctx).
		Delete(&models.MonitorProducer{}, producerID); gorm.IsFailResult(res) {

		return errors.WithMessage(errorByResult(res), op, "failed to delete monitor producer")
	}

	return nil
}

func (s *Storage) DeleteMonitor(
	ctx context.Context,
	monitorID int64,
) error {
	const op = "storage.mssql.monitor.DeleteMonitor"

	if res := s.db.WithContext(ctx).Delete(&models.Monitor{}, monitorID); gorm.IsFailResult(res) {
		return errors.WithMessage(errorByResult(res), op, "failed to delete monitor")
	}

	return nil
}
