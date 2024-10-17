package ssms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

func (s *Storage) SavePcType(
	ctx context.Context,
	name string,
	description string,
	processor *models.ProcessorData,
	videoCard *models.VideoCardData,
	monitor *models.MonitorData,
	ram *models.RamData,
) error {
	const op = "storage.ssms.pc.SavePc"
	stmt, args, err := squirrel.
		Expr(
			"EXEC InsertPcType ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
			name,
			description,
			processor.Producer,
			processor.Model,
			videoCard.Producer,
			videoCard.Model,
			monitor.Producer,
			monitor.Model,
			ram.Type,
			ram.Capacity,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, handleError(err))
	}

	stmt = replacePositionalParams(stmt, args)

	if _, err := s.db.ExecContext(ctx, stmt, args...); err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) SavePc(
	ctx context.Context,
	typeId int64,
	roomId int64,
	row int,
	place int,
) error {
	const op = "storage.ssms.pc.SavePc"

	stmt, args, err := squirrel.
		Insert("pc").
		Columns(
			"pc_type_id",
			"pc_room_id",
			"row",
			"place",
		).
		Values(
			typeId,
			roomId,
			row,
			place,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	if _, err := s.db.ExecContext(ctx, stmt, args...); err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeletePcType(
	ctx context.Context,
	typeId int64,
) error {
	const op = "storage.ssms.pc.DeletePc"

	stmt, args, err := squirrel.
		Delete("pc_types").
		Where(squirrel.Eq{
			"pc_type_id": typeId,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, handleError(err))
	}

	stmt = replacePositionalParams(stmt, args)

	if _, err := s.db.ExecContext(ctx, stmt, args...); err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) UpdatePcType(
	ctx context.Context,
	typeId int64,
	name string,
	description string,
	processor *models.ProcessorData,
	videoCard *models.VideoCardData,
	monitor *models.MonitorData,
	ram *models.RamData,
) error {
	const op = "storage.ssms.pc.UpdatePc"

	stmt, args, err := squirrel.Expr(
		"EXEC UpdatePcType ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
		typeId,
		name,
		description,
		processor.Producer,
		processor.Model,
		videoCard.Producer,
		videoCard.Model,
		monitor.Producer,
		monitor.Model,
		ram.Type,
		ram.Capacity,
	).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	var affectedCount int64
	if err := s.db.GetContext(ctx, &affectedCount, stmt, args...); err != nil {
		return fmt.Errorf("%s: failed to update pc type: %w", op, handleError(err))
	}

	if affectedCount == 0 {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}

	return nil
}