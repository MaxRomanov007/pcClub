package ssms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

func (s *Storage) ProcessorProducers(
	ctx context.Context,
) ([]models.ProcessorProducer, error) {
	const op = "storage.ssms.processor.ProcessorProducers"

	query, args, err := squirrel.Select("*").From("processor_producers").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var producers []models.ProcessorProducer
	if err := s.db.SelectContext(ctx, &producers, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to select processor producers: %w", op, handleError(err))
	}

	return producers, nil
}

func (s *Storage) Processors(
	ctx context.Context,
	producerId int64,
) ([]models.Processor, error) {
	const op = "storage.ssms.processor.Processors"

	query, args, err := squirrel.
		Select("*").
		From("processors").
		Where(squirrel.Eq{
			"processor_producer_id": producerId,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var processors []models.Processor
	if err := s.db.SelectContext(ctx, &processors, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to select processors: %w", op, handleError(err))
	}

	return processors, nil
}

func (s *Storage) SaveProcessorProducer(
	ctx context.Context,
	name string,
) error {
	const op = "storage.ssms.processor.SaveProcessorProducer"

	query, args, err := squirrel.Insert("processor_producers").Values(name).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to save processor producer: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) SaveProcessor(
	ctx context.Context,
	processor models.Processor,
) error {
	const op = "storage.ssms.processor.SaveProcessor"

	query, args, err := squirrel.
		Insert("processors").
		Columns(
			"processor_producer_id",
			"model",
		).
		Values(
			processor.ProcessorProducerID,
			processor.Model,
		).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to save processor: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeleteProcessorProducer(
	ctx context.Context,
	producerId int64,
) error {
	const op = "storage.ssms.processor.DeleteProcessorProducer"

	query, args, err := squirrel.
		Delete("processor_producers").
		Where(squirrel.Eq{
			"processor_producer_id": producerId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to delete processor producer: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeleteProcessor(
	ctx context.Context,
	processorId int64,
) error {
	const op = "storage.ssms.processor.DeleteProcessor"

	query, args, err := squirrel.
		Delete("processors").
		Where(squirrel.Eq{
			"processor_id": processorId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to delete processor: %w", op, handleError(err))
	}

	return nil
}
