package pcType

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/ssms"
)

func (s *Service) SavePcType(
	ctx context.Context,
	name string,
	description string,
	processor *models.ProcessorData,
	videoCard *models.VideoCardData,
	monitor *models.MonitorData,
	ram *models.RamData,
) error {
	const op = "services.pcClub.pc.SavePcType"

	err := s.owner.SavePcType(
		ctx,
		name,
		description,
		processor,
		videoCard,
		monitor,
		ram,
	)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
