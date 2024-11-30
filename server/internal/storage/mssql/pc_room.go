package mssql

import (
	"context"
	"server/internal/lib/api/database/gorm"
	"server/internal/lib/errors"
	"server/internal/models"
)

func (s *Storage) PcRoom(
	ctx context.Context,
	roomID int64,
) (models.PcRoom, error) {
	const op = "storage.mssql.pc_room.PcRoom"

	var room models.PcRoom
	if res := s.db.WithContext(ctx).First(&room, roomID); gorm.IsFailResult(res) {
		return models.PcRoom{}, errors.WithMessage(errorByResult(res), op, "failed to get pc room")
	}

	return room, nil
}

func (s *Storage) PcRooms(
	ctx context.Context,
	pcTypeId int64,
) ([]models.PcRoom, error) {
	const op = "storage.mssql.pc_room.PcRooms"

	sql := `
SELECT DISTINCT 
    dbo.pc_rooms.pc_room_id, 
    dbo.pc_rooms.name, 
    dbo.pc_rooms.places, 
    dbo.pc_rooms.rows, 
    dbo.pc_rooms.description
FROM pc_rooms
JOIN pc ON pc.pc_room_id = pc_rooms.pc_room_id
JOIN pc_types ON pc_types.pc_type_id = pc.pc_type_id
WHERE pc.pc_type_id = ?`

	var rooms []models.PcRoom
	if res := s.db.WithContext(ctx).Raw(sql, pcTypeId).Scan(&rooms); gorm.IsFailResult(res) {
		return nil, errors.WithMessage(errorByResult(res), op, "failed to get pc rooms")
	}

	return rooms, nil
}

func (s *Storage) SavePcRoom(
	ctx context.Context,
	room *models.PcRoom,
) (int64, error) {
	const op = "storage.mssql.pc_room.savePcRoom"

	if res := s.db.WithContext(ctx).Save(room); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save pc room")
	}

	return room.PcRoomID, nil
}

func (s *Storage) UpdatePcRoom(
	ctx context.Context,
	roomID int64,
	room *models.PcRoom,
) error {
	const op = "storage.mssql.pc_room.UpdatePcRoom"

	if res := s.db.WithContext(ctx).
		Where("pc_room_id = ?", roomID).
		Updates(&room); gorm.IsFailResult(res) {

		return errors.WithMessage(errorByResult(res), op, "failed to update pc room")
	}

	return nil
}

func (s *Storage) DeletePcRoom(
	ctx context.Context,
	roomID int64,
) error {
	const op = "storage.mssql.pc_room.DeletePcRoom"

	if res := s.db.WithContext(ctx).Delete(&models.PcRoom{}, roomID); gorm.IsFailResult(res) {
		return errors.WithMessage(errorByResult(res), op, "failed to delete pc room")
	}

	return nil
}
