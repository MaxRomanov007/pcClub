package ssms

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

type pcRoomNullable struct {
	PcRoomID    int64            `db:"pc_room_id"`
	Name        string           `db:"name"`
	Rows        int              `db:"rows"`
	Places      int              `db:"places"`
	Description sql.Null[string] `db:"description"`
}

func (null *pcRoomNullable) parse() models.PcRoom {
	return models.PcRoom{
		PcRoomID:    null.PcRoomID,
		Name:        null.Name,
		Rows:        null.Rows,
		Places:      null.Places,
		Description: null.Description.V,
	}
}

func (s *Storage) PcRoom(
	ctx context.Context,
	roomId int64,
) (models.PcRoom, error) {
	const op = "storage.ssms.pc-room.PcRoom"

	query, args, err := squirrel.
		Select("*").
		From("pc_rooms").
		Where(squirrel.Eq{
			"pc_room_id": roomId,
		}).ToSql()
	if err != nil {
		return models.PcRoom{}, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var pcRoom pcRoomNullable
	if err := s.db.GetContext(ctx, &pcRoom, query, args...); err != nil {
		return models.PcRoom{}, handleError(err)
	}

	return pcRoom.parse(), nil
}

func (s *Storage) SavePcRoom(
	ctx context.Context,
	pcRoom models.PcRoom,
) error {
	const op = "storage.ssms.pc-room.savePcRoom"

	query, args, err := squirrel.Insert("pc_rooms").
		Columns(
			"name",
			"rows",
			"places",
			"description",
		).
		Values(
			pcRoom.Name,
			pcRoom.Rows,
			pcRoom.Places,
			nullString(pcRoom.Description),
		).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) UpdatePcRoom(
	ctx context.Context,
	pcRoom models.PcRoom,
) error {
	const op = "storage.ssms.pc-room.UpdatePcRoom"

	stmt := squirrel.Update("pc_rooms").
		Where(squirrel.Eq{
			"pc_room_id": pcRoom.PcRoomID,
		})
	stmt = setIfNotZero(stmt, "name", pcRoom.Name)
	stmt = setIfNotZero(stmt, "rows", pcRoom.Rows)
	stmt = setIfNotZero(stmt, "places", pcRoom.Places)
	stmt = setIfNotZeroNullString(stmt, "description", pcRoom.Description)

	query, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeletePcRoom(
	ctx context.Context,
	roomId int64,
) error {
	const op = "storage.ssms.pc-room.DeletePcRoom"

	query, args, err := squirrel.
		Delete("pc_rooms").
		Where(squirrel.Eq{
			"pc_room_id": roomId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build statement: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to exec query: %w", op, handleError(err))
	}
	return nil
}
