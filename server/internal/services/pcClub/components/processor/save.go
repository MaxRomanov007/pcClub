package processor

import (
	"context"
	"server/internal/lib/errors"
	"server/internal/models"
	"server/internal/services/pcClub/components"
)

func (s *Service) SaveProcessorProducer(
	ctx context.Context,
	producer *models.ProcessorProducer,
) (int64, error) {
	const op = "services.pcClub.components.processor.SaveProcessorProducer"

	id, err := s.owner.SaveProcessorProducer(ctx, producer)
	if err != nil {
		return 0, errors.WithMessage(components.HandleStorageError(err), op, "failed to save processor producer in mssql")
	}

	return id, nil
}

func (s *Service) SaveProcessor(
	ctx context.Context,
	processor *models.Processor,
) (int64, error) {
	const op = "services.pcClub.components.processor.SaveProcessor"

	id, err := s.owner.SaveProcessor(ctx, processor)
	if err != nil {
		return 0, errors.WithMessage(components.HandleStorageError(err), op, "failed to save processor in mssql")
	}

	return id, nil
}
