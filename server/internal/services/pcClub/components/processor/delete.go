package processor

import (
	"context"
	"errors"
	"fmt"
	"server/internal/services/pcClub/components"
	"server/internal/storage/ssms"
)

func (s *Service) DeleteProcessorProducer(
	ctx context.Context,
	producerId int64,
) error {
	const op = "services.pcClub.components.processor.DeleteProcessorProducer"

	if err := s.owner.DeleteProcessorProducer(ctx, producerId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete processor producer: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteProcessor(
	ctx context.Context,
	processorId int64,
) error {
	const op = "services.pcClub.components.processor.DeleteProcessor"

	if err := s.owner.DeleteProcessor(ctx, processorId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete processor: %w", op, err)
	}

	return nil
}
