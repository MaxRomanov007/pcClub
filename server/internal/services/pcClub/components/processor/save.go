package processor

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/ssms"
)

func (s *Service) SaveProcessorProducer(
	ctx context.Context,
	name string,
) error {
	const op = "services.pcClub.components.processor.SaveProcessorProducer"

	if err := s.owner.SaveProcessorProducer(ctx, name); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save processor producer: %w", op, err)
	}

	return nil
}

func (s *Service) SaveProcessor(
	ctx context.Context,
	processor models.Processor,
) error {
	const op = "services.pcClub.components.processor.SaveProcessor"

	if err := s.owner.SaveProcessor(ctx, processor); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save processor: %w", op, err)
	}

	return nil
}
