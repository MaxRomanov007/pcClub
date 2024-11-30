package mssql

import (
	"context"
	"server/internal/lib/api/database/gorm"
	"server/internal/lib/errors"
	"server/internal/models"
)

func (s *Storage) Dishes(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Dish, error) {
	const op = "storage.mssql.dish.Dishes"

	var dishes []models.Dish
	if res := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&dishes); gorm.IsFailResult(res) {

		return nil, errors.WithMessage(errorByResult(res), op, "failed to get dishes")
	}

	return dishes, nil
}

func (s *Storage) Dish(
	ctx context.Context,
	dishID int64,
) (models.Dish, error) {
	const op = "storage.mssql.dish.Dish"

	var dish models.Dish
	if res := s.db.WithContext(ctx).First(&dish, dishID); gorm.IsFailResult(res) {
		return models.Dish{}, errors.WithMessage(errorByResult(res), op, "failed to get dish")
	}

	return dish, nil
}

func (s *Storage) SaveDish(
	ctx context.Context,
	dish *models.Dish,
) (int64, error) {
	const op = "storage.mssql.dish.SaveDish"

	if res := s.db.WithContext(ctx).Save(dish); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save dish")
	}

	return dish.DishID, nil
}

func (s *Storage) UpdateDish(
	ctx context.Context,
	dishID int64,
	dish *models.Dish,
) error {
	const op = "storage.mssql.dish.UpdateDish"

	if res := s.db.WithContext(ctx).
		Where("dish_id = ?", dishID).
		Updates(&dish); gorm.IsFailResult(res) {

		return errors.WithMessage(errorByResult(res), op, "failed to update dish")
	}

	return nil
}

func (s *Storage) DeleteDish(
	ctx context.Context,
	dishId int64,
) error {
	const op = "storage.mssql.dish.DeleteDish"

	if res := s.db.WithContext(ctx).Delete(&models.Dish{}, dishId); gorm.IsFailResult(res) {
		return errors.WithMessage(errorByResult(res), op, "failed to delete dish")
	}

	return nil
}
