package videoCard

import (
	"context"
	"server/internal/lib/errors"
)

func (s *Service) DeleteVideoCardProducer(
	ctx context.Context,
	producerID int64,
) error {
	const op = "services.pcClub.components.videoCard.DeleteVideoCardProducer"

	if err := s.owner.DeleteVideoCardProducer(ctx, producerID); err != nil {
		return errors.WithMessage(err, op, "failed to delete video card producer from mssql")
	}

	return nil
}

func (s *Service) DeleteVideoCard(
	ctx context.Context,
	videoCardID int64,
) error {
	const op = "services.pcClub.components.videoCard.DeleteVideoCard"

	if err := s.owner.DeleteVideoCard(ctx, videoCardID); err != nil {
		return errors.WithMessage(err, op, "failed to delete video card from mssql")
	}

	return nil
}
