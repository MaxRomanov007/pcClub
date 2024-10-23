package pcType

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/ssms"
)

func (s *Service) UpdatePcType(
	ctx context.Context,
	typeId int64,
	name string,
	description string,
	processor *models.ProcessorData,
	videoCard *models.VideoCardData,
	monitor *models.MonitorData,
	ram *models.RamData,
) error {
	const op = "services.pcClub.pc.UpdatePcType"

	if err := s.owner.UpdatePcType(
		ctx,
		typeId,
		name,
		description,
		processor,
		videoCard,
		monitor,
		ram,
	); err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return fmt.Errorf("%s: failed to update pc type: %w", op, err)
	}

	return nil
}
