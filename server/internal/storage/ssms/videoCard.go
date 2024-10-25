package ssms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

func (s *Storage) VideoCardProducers(
	ctx context.Context,
) ([]models.VideoCardProducer, error) {
	const op = "storage.ssms.videoCard.VideoCardProducers"

	query, args, err := squirrel.Select("*").From("video_card_producers").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var producers []models.VideoCardProducer
	if err := s.db.SelectContext(ctx, &producers, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to select video card producers: %w", op, handleError(err))
	}

	return producers, nil
}

func (s *Storage) VideoCards(
	ctx context.Context,
	producerId int64,
) ([]models.VideoCard, error) {
	const op = "storage.ssms.videoCard.VideoCards"

	query, args, err := squirrel.
		Select("*").
		From("video_cards").
		Where(squirrel.Eq{
			"video_card_producer_id": producerId,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var videoCards []models.VideoCard
	if err := s.db.SelectContext(ctx, &videoCards, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to select video cards: %w", op, handleError(err))
	}

	return videoCards, nil
}

func (s *Storage) SaveVideoCardProducer(
	ctx context.Context,
	name string,
) error {
	const op = "storage.ssms.videoCard.SaveVideoCardProducer"

	query, args, err := squirrel.Insert("video_card_producers").Values(name).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to save video card producer: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) SaveVideoCard(
	ctx context.Context,
	videoCard models.VideoCard,
) error {
	const op = "storage.ssms.videoCard.SaveVideoCard"

	query, args, err := squirrel.
		Insert("video_cards").
		Columns(
			"video_card_producer_id",
			"model",
		).
		Values(
			videoCard.VideoCardProducerID,
			videoCard.Model,
		).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to save video card: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeleteVideoCardProducer(
	ctx context.Context,
	producerId int64,
) error {
	const op = "storage.ssms.videoCard.DeleteVideoCardProducer"

	query, args, err := squirrel.
		Delete("video_card_producers").
		Where(squirrel.Eq{
			"video_card_producer_id": producerId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to delete video card producer: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeleteVideoCard(
	ctx context.Context,
	videoCardId int64,
) error {
	const op = "storage.ssms.videoCard.DeleteVideoCard"

	query, args, err := squirrel.
		Delete("video_cards").
		Where(squirrel.Eq{
			"video_card_id": videoCardId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to delete video card: %w", op, handleError(err))
	}

	return nil
}
