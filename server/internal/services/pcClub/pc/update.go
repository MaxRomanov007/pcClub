package pc

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
		if errors.Is(err, ssms.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return fmt.Errorf("%s: failed to update pc type: %w", op, err)
	}

	return nil
}

func (s *Service) UpdatePc(
	ctx context.Context,
	pcId int64,
	typeId int64,
	roomId int64,
	statusId int64,
	row int,
	place int,
	description string,
) error {
	const op = "services.pcClub.pc.UpdatePc"

	if err := s.owner.UpdatePc(
		ctx,
		pcId,
		typeId,
		roomId,
		statusId,
		row,
		place,
		description,
	); err != nil {
		if errors.Is(err, ssms.ErrCheckFailed) {
			return fmt.Errorf("%s: %w", op, ErrConstraint)
		}
		if errors.Is(err, ssms.ErrAlreadyExists) {
			return fmt.Errorf("%s: %w", op, ErrAlreadyExists)
		}
		if errors.Is(err, ssms.ErrReferenceNotExists) {
			return fmt.Errorf("%s: %w", op, ErrReferenceNotExists)
		}
		return fmt.Errorf("%s: failed to update pc: %w", op, err)
	}

	return nil
}
