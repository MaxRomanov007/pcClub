package processor

import (
	"context"
	"server/internal/lib/errors"
)

func (s *Service) DeleteProcessorProducer(
	ctx context.Context,
	producerID int64,
) error {
	const op = "services.pcClub.components.processor.DeleteProcessorProducer"

	if err := s.owner.DeleteProcessorProducer(ctx, producerID); err != nil {
		return errors.WithMessage(err, op, "failed to delete processor producer from mssql")
	}

	return nil
}

func (s *Service) DeleteProcessor(
	ctx context.Context,
	processorID int64,
) error {
	const op = "services.pcClub.components.processor.DeleteProcessor"

	if err := s.owner.DeleteProcessor(ctx, processorID); err != nil {
		return errors.WithMessage(err, op, "failed to delete processor from mssql")
	}

	return nil
}
