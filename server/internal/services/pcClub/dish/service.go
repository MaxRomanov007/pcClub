package dish

import (
	"context"
	"server/internal/models"
)

type provider interface {
	Dishes(
		ctx context.Context,
		limit int64,
		offset int64,
	) (dishes []models.DishData, err error)

	Dish(
		ctx context.Context,
		dishId int64,
	) (dish models.DishData, err error)
}

type owner interface {
	SaveDish(
		ctx context.Context,
		dish models.DishData,
	) (err error)

	UpdateDish(
		ctx context.Context,
		dish models.DishData,
	) (err error)

	DeleteDish(
		ctx context.Context,
		dishId int64,
	) (err error)
}

type redisProvider interface {
	Value(
		ctx context.Context,
		key string,
		value interface{},
	) (err error)
}

type redisOwner interface {
	Set(
		ctx context.Context,
		key string,
		value interface{},
	) (err error)
}

type Service struct {
	provider      provider
	owner         owner
	redisProvider redisProvider
	redisOwner    redisOwner
}

func New(provider provider, owner owner, redisProvider redisProvider, redisOwner redisOwner) *Service {
	return &Service{
		provider:      provider,
		owner:         owner,
		redisProvider: redisProvider,
		redisOwner:    redisOwner,
	}
}
