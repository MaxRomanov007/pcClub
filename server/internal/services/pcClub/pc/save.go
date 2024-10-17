package pc

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/sqlServer"
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
	if errors.Is(err, sqlServer.ErrAlreadyExists) {
		return fmt.Errorf("%s: %w", op, ErrAlreadyExists)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) SavePc(
	ctx context.Context,
	typeId int64,
	roomId int64,
	row int,
	place int,
) error {
	const op = "services.pcClub.pc.SavePc"

	err := s.owner.SavePc(
		ctx,
		typeId,
		roomId,
		row,
		place,
	)
	if errors.Is(err, sqlServer.ErrAlreadyExists) {
		return fmt.Errorf("%s: %w", op, ErrAlreadyExists)
	}
	if err != nil {
		return fmt.Errorf("%s: failed to save pc: %w", op, err)
	}

	return nil
}
