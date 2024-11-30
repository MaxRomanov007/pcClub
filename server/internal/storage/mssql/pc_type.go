package mssql

import (
	"context"
	"server/internal/lib/api/database/gorm"
	"server/internal/lib/errors"
	"server/internal/models"
)

func (s *Storage) PcTypes(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.PcType, error) {
	const op = "storage.mssql.pc_type.PcTypes"

	var types []models.PcType
	if res := s.db.WithContext(ctx).
		Limit(limit).Offset(offset).
		Find(&types); gorm.IsFailResult(res) {

		return nil, errors.WithMessage(errorByResult(res), op, "failed to get pc types")
	}

	return types, nil
}

func (s *Storage) PcType(
	ctx context.Context,
	typeID int64,
) (models.PcType, error) {
	const op = "storage.mssql.pc_type.PcType"

	var pcType models.PcType
	if res := s.db.WithContext(ctx).
		Preload("Processor").Preload("Processor.ProcessorProducer").
		Preload("VideoCard").Preload("VideoCard.VideoCardProducer").
		Preload("Monitor").Preload("Monitor.MonitorProducer").
		Preload("RAM").Preload("RAM.RAMType").
		First(&pcType, typeID); gorm.IsFailResult(res) {

		return models.PcType{}, errors.WithMessage(errorByResult(res), op, "failed to get pc type")
	}

	return pcType, nil
}

func (s *Storage) SavePcType(
	ctx context.Context,
	pcType *models.PcType,
) (int64, error) {
	const op = "storage.mssql.pc_type.SavePc"

	if res := s.db.WithContext(ctx).Save(pcType); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save pc type")
	}

	return pcType.PcTypeID, nil
}

func (s *Storage) UpdatePcType(
	ctx context.Context,
	typeID int64,
	pcType *models.PcType,
) error {
	const op = "storage.mssql.pc_type.UpdatePc"

	if res := s.db.WithContext(ctx).
		Where("pc_type_id = ?", typeID).
		Updates(pcType); gorm.IsFailResult(res) {

		return errors.WithMessage(errorByResult(res), op, "failed to update pc type")
	}

	return nil
}

func (s *Storage) DeletePcType(
	ctx context.Context,
	typeID int64,
) error {
	const op = "storage.mssql.pc_type.DeletePc"

	if res := s.db.WithContext(ctx).Delete(&models.PcType{}, typeID); gorm.IsFailResult(res) {
		return errors.WithMessage(errorByResult(res), op, "failed to delete pc type")
	}

	return nil
}
