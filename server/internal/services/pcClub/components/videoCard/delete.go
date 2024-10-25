package videoCard

import (
	"context"
	"errors"
	"fmt"
	"server/internal/services/pcClub/components"
	"server/internal/storage/ssms"
)

func (s *Service) DeleteVideoCardProducer(
	ctx context.Context,
	producerId int64,
) error {
	const op = "services.pcClub.components.videoCard.DeleteVideoCardProducer"

	if err := s.owner.DeleteVideoCardProducer(ctx, producerId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete video card producer: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteVideoCard(
	ctx context.Context,
	videoCardId int64,
) error {
	const op = "services.pcClub.components.videoCard.DeleteVideoCard"

	if err := s.owner.DeleteVideoCard(ctx, videoCardId); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, components.HandleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to delete video card: %w", op, err)
	}

	return nil
}
