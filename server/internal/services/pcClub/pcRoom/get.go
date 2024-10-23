package pcRoom

import (
	"context"
	"errors"
	"fmt"
	"server/internal/models"
	"server/internal/storage/redis"
	"server/internal/storage/ssms"
)

func (s *Service) PcRoom(
	ctx context.Context,
	pcRoomId int64,
) (models.PcRoom, error) {
	const op = "services.pcClub.pcRoom.PcRoom"

	var pcRoom models.PcRoom
	err := s.redisProvider.Value(
		ctx,
		fmt.Sprintf("pc_room:%d", pcRoomId),
		&pcRoom,
	)
	if err == nil {
		return pcRoom, nil
	}
	if !errors.Is(err, redis.ErrNotFound) {
		return models.PcRoom{}, fmt.Errorf("%s: failed to get pc room from redis: %w", op, err)
	}

	pcRoom, err = s.provider.PcRoom(ctx, pcRoomId)
	if err != nil {
		var ssmsErr *ssms.Error
		if errors.As(err, &ssmsErr) {
			return models.PcRoom{}, fmt.Errorf("%s: %w", op, handleStorageError(ssmsErr))
		}
		return models.PcRoom{}, fmt.Errorf("%s: failed to get pc room from storage: %w", op, err)
	}

	if err := s.redisOwner.Set(
		ctx,
		fmt.Sprintf("pc_room:%d", pcRoomId),
		pcRoom,
	); err != nil {
		return models.PcRoom{}, fmt.Errorf("%s: failed to send pc room in redis: %w", op, err)
	}

	return pcRoom, nil
}
