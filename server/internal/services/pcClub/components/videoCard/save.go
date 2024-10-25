package videoCard

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/services/pcClub/components"
	"server/internal/storage/ssms"
)

func (s *Service) SaveVideoCardProducer(
	ctx context.Context,
	name string,
) error {
	const op = "services.pcClub.components.videoCard.SaveVideoCardProducer"

	if err := s.owner.SaveVideoCardProducer(ctx, name); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save video card producer: %w", op, err)
	}

	return nil
}

func (s *Service) SaveVideoCard(
	ctx context.Context,
	videoCard models.VideoCard,
) error {
	const op = "services.pcClub.components.videoCard.SaveVideoCard"

	if err := s.owner.SaveVideoCard(ctx, videoCard); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to save video card: %w", op, err)
	}

	return nil
}
