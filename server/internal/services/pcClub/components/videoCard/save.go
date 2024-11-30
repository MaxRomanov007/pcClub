package videoCard

import (
	"context"
	"server/internal/lib/errors"
	"server/internal/models"
	"server/internal/services/pcClub/components"
)

func (s *Service) SaveVideoCardProducer(
	ctx context.Context,
	producer *models.VideoCardProducer,
) (int64, error) {
	const op = "services.pcClub.components.videoCard.SaveVideoCardProducer"

	id, err := s.owner.SaveVideoCardProducer(ctx, producer)
	if err != nil {
		return 0, errors.WithMessage(components.HandleStorageError(err), op, "failed to save video card producer in mssql")
	}

	return id, nil
}

func (s *Service) SaveVideoCard(
	ctx context.Context,
	videoCard *models.VideoCard,
) (int64, error) {
	const op = "services.pcClub.components.videoCard.SaveVideoCard"

	id, err := s.owner.SaveVideoCard(ctx, videoCard)
	if err != nil {
		return 0, errors.WithMessage(components.HandleStorageError(err), op, "failed to save video card in mssql")
	}

	return id, nil
}
