package ssms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

func (s *Storage) Pcs(
	ctx context.Context,
	typeId int64,
	isAvailable bool,
) ([]models.PcData, error) {
	const op = "storage.ssms.pc.Pc"

	squirrelStmt := squirrel.Select(
		"pc.pc_id",
		"pc.row",
		"pc.place",
		"pc.description",
		"pc.pc_room_id",
		"pc_statuses.name AS status",
	).
		From("pc").
		Join("pc_statuses ON pc.pc_status_id = pc_statuses.pc_status_id").
		Where(squirrel.Eq{
			"pc.pc_type_id": typeId,
		})

	if isAvailable {
		squirrelStmt = squirrelStmt.
			Where(squirrel.Eq{
				"pc_statuses.name": "available",
			})
	}

	stmt, args, err := squirrelStmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	var pcs []models.Pc
	if err := s.db.SelectContext(ctx, &pcs, stmt, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to get pcs: %w", op, handleError(err))
	}

	pcsData := make([]models.PcData, len(pcs))
	for i, pc := range pcs {
		pcsData[i] = models.PcData{
			PcID:        pc.PcID,
			Row:         pc.Row,
			Place:       pc.Place,
			Description: pc.Description.V,
			PcRoomID:    pc.PcRoomID,
			Status:      pc.Status,
		}
	}

	return pcsData, nil
}

func (s *Storage) SavePc(
	ctx context.Context,
	typeId int64,
	roomId int64,
	row int,
	place int,
	description string,
) error {
	const op = "storage.ssms.pc.SavePc"

	stmt, args, err := squirrel.
		Insert("pc").
		Columns(
			"pc_type_id",
			"pc_room_id",
			"row",
			"place",
			"[description]",
		).
		Values(
			typeId,
			roomId,
			row,
			place,
			nullString(description),
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

func (s *Storage) UpdatePc(
	ctx context.Context,
	pcId int64,
	typeId int64,
	roomId int64,
	statusId int64,
	row int,
	place int,
	description string,
) error {
	const op = "storage.ssms.pc.UpdatePc"

	stmt := squirrel.
		Update("pc").
		Where(squirrel.Eq{
			"pc_id": pcId,
		})

	stmt = setIfNotZero(stmt, "pc_type_id", typeId)
	stmt = setIfNotZero(stmt, "pc_room_id", roomId)
	stmt = setIfNotZero(stmt, "pc_status_id", statusId)
	stmt = setIfNotZero(stmt, "row", row)
	stmt = setIfNotZero(stmt, "place", place)
	stmt = setIfNotZeroNullString(stmt, "description", description)

	query, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeletePc(
	ctx context.Context,
	pcId int64,
) error {
	const op = "storage.ssms.pc.DeletePc"

	query, args, err := squirrel.Delete("pc").
		Where(squirrel.Eq{
			"pc_id": pcId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}
