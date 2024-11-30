package mssql

import (
	"golang.org/x/net/context"
	"server/internal/lib/api/database/gorm"
	"server/internal/lib/errors"
	"server/internal/models"
)

func (s *Storage) ProcessorProducers(
	ctx context.Context,
) ([]models.ProcessorProducer, error) {
	const op = "storage.mssql.processor.ProcessorProducers"

	var producers []models.ProcessorProducer
	if res := s.db.WithContext(ctx).Find(&producers); gorm.IsFailResult(res) {
		return nil, errors.WithMessage(errorByResult(res), op, "failed to get processor producers")
	}

	return producers, nil
}

func (s *Storage) Processors(
	ctx context.Context,
	producerID int64,
) ([]models.Processor, error) {
	const op = "storage.mssql.monitor.MonitorProducer"

	var processors []models.Processor
	if res := s.db.WithContext(ctx).
		Where("processor_producer_id = ?", producerID).
		Find(&processors); gorm.IsFailResult(res) {

		return nil, errors.WithMessage(
			errorByResult(res),
			op, "failed to get processors",
		)
	}

	return processors, nil
}

func (s *Storage) SaveProcessorProducer(
	ctx context.Context,
	producer *models.ProcessorProducer,
) (int64, error) {
	const op = "storage.mssql.processor.SaveProcessorProducer"

	if res := s.db.WithContext(ctx).Save(producer); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save processor producer")
	}

	return producer.ProcessorProducerID, nil
}

func (s *Storage) SaveProcessor(
	ctx context.Context,
	processor *models.Processor,
) (int64, error) {
	const op = "storage.mssql.processor.SaveProcessor"

	if res := s.db.WithContext(ctx).Save(processor); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save processor")
	}

	return processor.ProcessorID, nil
}

func (s *Storage) DeleteProcessorProducer(
	ctx context.Context,
	producerID int64,
) error {
	const op = "storage.mssql.processor.DeleteProcessorProducer"

	if res := s.db.WithContext(ctx).
		Delete(&models.ProcessorProducer{}, producerID); gorm.IsFailResult(res) {

		return errors.WithMessage(errorByResult(res), op, "failed to delete processor producer")
	}

	return nil
}

func (s *Storage) DeleteProcessor(
	ctx context.Context,
	processorID int64,
) error {
	const op = "storage.mssql.processor.DeleteProcessor"

	if res := s.db.WithContext(ctx).Delete(&models.Processor{}, processorID); gorm.IsFailResult(res) {
		return errors.WithMessage(errorByResult(res), op, "failed to delete processor")
	}

	return nil
}
