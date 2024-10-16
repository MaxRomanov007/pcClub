package user

import (
	"context"
	"errors"
	"server/internal/config"
	"server/internal/models"
)

type provider interface {
	User(
		ctx context.Context,
		uid int64,
	) (user models.UserData, err error)

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
		user models.User,
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

var (
	ErrInvalidCredentials = errors.New("credentials are not valid")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrAccessDenied       = errors.New("access denied")
	ErrUserNotFound       = errors.New("user not found")
)

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