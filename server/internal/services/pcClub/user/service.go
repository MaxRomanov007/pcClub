package user

import (
	"context"
	"server/internal/config"
	"server/internal/models"
)

type provider interface {
	User(
		ctx context.Context,
		uid int64,
	) (user models.User, err error)

	UserByEmail(
		ctx context.Context,
		email string,
	) (user models.User, err error)

	UserRole(
		ctx context.Context,
		uid int64,
	) (role string, err error)
}

type owner interface {
	SaveUser(
		ctx context.Context,
		user *models.User,
	) (id int64, err error)

	DeleteUser(
		ctx context.Context,
		uid int64,
	) (err error)
}

type Service struct {
	cfg          *config.UserConfig
	userProvider provider
	userOwner    owner
}

func New(
	cfg *config.UserConfig,
	userProvider provider,
	userOwner owner,
) *Service {
	return &Service{
		cfg:          cfg,
		userProvider: userProvider,
		userOwner:    userOwner,
	}
}
