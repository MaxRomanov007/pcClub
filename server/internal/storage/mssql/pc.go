package mssql

import (
	"context"
	"server/internal/lib/api/database/gorm"
	"server/internal/lib/errors"
	"server/internal/models"
)

func (s *Storage) Pcs(
	ctx context.Context,
	typeID int64,
	isAvailable bool,
) ([]models.Pc, error) {
	const op = "storage.mssql.pc.Pc"

	var pcs []models.Pc
	db := s.db.WithContext(ctx).Where("pc_type_id = ?", typeID)
	if isAvailable {
		db = db.Joins("PcStatus").Where("PcStatus.Name = ?", AvailablePcStatus)
	}
	if res := db.Find(&pcs); gorm.IsFailResult(res) {
		return nil, errors.WithMessage(errorByResult(res), op, "failed to get pcs")
	}

	return pcs, nil
}

func (s *Storage) SavePc(
	ctx context.Context,
	pc *models.Pc,
) (int64, error) {
	const op = "storage.mssql.pc.SavePc"

	if res := s.db.WithContext(ctx).Save(pc); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save pc")
	}

	return pc.PcID, nil
}

func (s *Storage) UpdatePc(
	ctx context.Context,
	pcID int64,
	pc *models.Pc,
) error {
	const op = "storage.mssql.pc.UpdatePc"

	if res := s.db.WithContext(ctx).
		Where("pc_id = ?", pcID).
		Updates(pc); gorm.IsFailResult(res) {

		return errors.WithMessage(errorByResult(res), op, "failed to update pc")
	}

	return nil
}

func (s *Storage) DeletePc(
	ctx context.Context,
	pcID int64,
) error {
	const op = "storage.mssql.pc.DeletePc"

	if res := s.db.WithContext(ctx).Delete(&models.Pc{}, pcID); gorm.IsFailResult(res) {
		return errors.WithMessage(errorByResult(res), op, "failed to delete pc")
	}

	return nil
}
