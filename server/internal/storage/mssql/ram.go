package mssql

import (
	"context"
	"server/internal/lib/api/database/gorm"
	"server/internal/lib/errors"
	"server/internal/models"
)

func (s *Storage) RamTypes(
	ctx context.Context,
) ([]models.RAMType, error) {
	const op = "storage.mssql.ram.RamTypes"

	var ramTypes []models.RAMType
	if res := s.db.WithContext(ctx).Find(&ramTypes); gorm.IsFailResult(res) {
		return nil, errors.WithMessage(errorByResult(res), op, "failed to get ram types")
	}

	return ramTypes, nil
}

func (s *Storage) Rams(
	ctx context.Context,
	ramTypeID int64,
) ([]models.RAM, error) {
	const op = "storage.mssql.ram.RamType"

	var rams []models.RAM
	if res := s.db.WithContext(ctx).
		Where("ram_type_id = ?", ramTypeID).
		Find(&rams, ramTypeID); gorm.IsFailResult(res) {

		return nil, errors.WithMessage(
			errorByResult(res),
			op, "failed to get rams",
		)
	}

	return rams, nil
}

func (s *Storage) SaveRamType(
	ctx context.Context,
	ramType *models.RAMType,
) (int64, error) {
	const op = "storage.mssql.ram.SaveRamType"

	if res := s.db.WithContext(ctx).Save(ramType); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save ram type")
	}

	return ramType.RAMTypeID, nil
}

func (s *Storage) SaveRam(
	ctx context.Context,
	ram *models.RAM,
) (int64, error) {
	const op = "storage.mssql.ram.SaveRam"

	if res := s.db.WithContext(ctx).Save(ram); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save ram")
	}

	return ram.RAMID, nil
}

func (s *Storage) DeleteRamType(
	ctx context.Context,
	ramTypeID int64,
) error {
	const op = "storage.mssql.ram.DeleteRamType"

	if res := s.db.WithContext(ctx).
		Delete(&models.RAMType{}, ramTypeID); gorm.IsFailResult(res) {

		return errors.WithMessage(errorByResult(res), op, "failed to delete ram type")
	}

	return nil
}

func (s *Storage) DeleteRam(
	ctx context.Context,
	ramID int64,
) error {
	const op = "storage.mssql.ram.DeleteRam"

	if res := s.db.WithContext(ctx).Delete(&models.RAM{}, ramID); gorm.IsFailResult(res) {
		return errors.WithMessage(errorByResult(res), op, "failed to delete ram")
	}

	return nil
}
