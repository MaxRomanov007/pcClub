package mssql

import (
	"context"
	"server/internal/lib/api/database/gorm"
	"server/internal/lib/errors"
	"server/internal/models"
)

func (s *Storage) VideoCardProducers(
	ctx context.Context,
) ([]models.VideoCardProducer, error) {
	const op = "storage.mssql.video_card.VideoCardProducers"

	var producers []models.VideoCardProducer
	if res := s.db.WithContext(ctx).Find(&producers); gorm.IsFailResult(res) {
		return nil, errors.WithMessage(errorByResult(res), op, "failed to get video card producers")
	}

	return producers, nil
}

func (s *Storage) VideoCards(
	ctx context.Context,
	producerID int64,
) ([]models.VideoCard, error) {
	const op = "storage.mssql.video_card.VideoCardProducer"

	var cards []models.VideoCard
	if res := s.db.WithContext(ctx).
		Where("video_card_producer = ?", producerID).
		Find(&cards, producerID); gorm.IsFailResult(res) {

		return nil, errors.WithMessage(
			errorByResult(res),
			op, "failed to get video cards",
		)
	}

	return cards, nil
}

func (s *Storage) SaveVideoCardProducer(
	ctx context.Context,
	producer *models.VideoCardProducer,
) (int64, error) {
	const op = "storage.mssql.video_card.SaveVideoCardProducer"

	if res := s.db.WithContext(ctx).Save(producer); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save video card producer")
	}

	return producer.VideoCardProducerID, nil
}

func (s *Storage) SaveVideoCard(
	ctx context.Context,
	card *models.VideoCard,
) (int64, error) {
	const op = "storage.mssql.video_card.SaveVideoCard"

	if res := s.db.WithContext(ctx).Save(card); gorm.IsFailResult(res) {
		return 0, errors.WithMessage(errorByResult(res), op, "failed to save video card")
	}

	return card.VideoCardID, nil
}

func (s *Storage) DeleteVideoCardProducer(
	ctx context.Context,
	producerID int64,
) error {
	const op = "storage.mssql.video_card.DeleteVideoCardProducer"

	if res := s.db.WithContext(ctx).
		Delete(&models.VideoCardProducer{}, producerID); gorm.IsFailResult(res) {

		return errors.WithMessage(errorByResult(res), op, "failed to delete video card producer")
	}

	return nil
}

func (s *Storage) DeleteVideoCard(
	ctx context.Context,
	cardID int64,
) error {
	const op = "storage.mssql.video_card.DeleteVideoCard"

	if res := s.db.WithContext(ctx).Delete(&models.VideoCard{}, cardID); gorm.IsFailResult(res) {
		return errors.WithMessage(errorByResult(res), op, "failed to delete video card")
	}

	return nil
}
