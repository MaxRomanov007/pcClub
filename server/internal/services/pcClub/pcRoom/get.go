package pcRoom

import (
	"context"
	"errors"
	"fmt"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
	"server/internal/storage/redis"
)

func (s *Service) PcRoom(
	ctx context.Context,
	roomID int64,
) (models.PcRoom, error) {
	const op = "services.pcClub.pcRoom.PcRoom"

	var room models.PcRoom
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("pc_room:%d", roomID),
		&room,
	)
	if err == nil {
		return room, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return models.PcRoom{}, errors2.WithMessage(err, op, "failed to get pc room from redis")
	}

	room, err = s.provider.PcRoom(ctx, roomID)
	if err != nil {
		return models.PcRoom{}, errors2.WithMessage(HandleStorageError(err), op, "failed to get pc room from mssql")
	}

	if err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("pc_room:%d", roomID),
		room,
	); err != nil {
		return models.PcRoom{}, errors2.WithMessage(err, op, "failed to save pc room in redis")
	}

	return room, nil
}

func (s *Service) PcRooms(
	ctx context.Context,
	pcTypeID int64,
) ([]models.PcRoom, error) {
	const op = "services.pcClub.pcRoom.PcRooms"

	var rooms []models.PcRoom
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("pc_rooms:%d", pcTypeID),
		&rooms,
	)
	if err == nil {
		return rooms, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return nil, errors2.WithMessage(err, op, "failed to get pc rooms from redis")
	}

	rooms, err = s.provider.PcRooms(ctx, pcTypeID)
	if err != nil {
		return nil, errors2.WithMessage(HandleStorageError(err), op, "failed to get pc rooms from mssql")
	}

	if err = s.redisOwner.Set(
		ctx,
		fmt.Sprintf("pc_room:%d", pcTypeID),
		rooms,
	); err != nil {
		return nil, errors2.WithMessage(err, op, "failed to save pc rooms in redis")
	}

	return rooms, nil
}
