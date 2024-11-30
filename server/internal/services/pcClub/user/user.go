package user

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	errors2 "server/internal/lib/errors"
	"server/internal/models"
	"server/internal/storage/mssql"
)

func (s *Service) User(
	ctx context.Context,
	uid int64,
) (models.User, error) {
	const op = "services.pcClub.user.User"

	user, err := s.userProvider.User(ctx, uid)
	if err != nil {
		return models.User{}, errors2.WithMessage(HandleStorageError(err), op, "failed to get user from mssql")
	}

	return user, nil
}

func (s *Service) UserByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "services.pcClub.user.UserByEmail"

	user, err := s.userProvider.UserByEmail(ctx, email)
	if err != nil {
		return models.User{}, errors2.WithMessage(HandleStorageError(err), op, "failed to get user by email from mssql")
	}

	return user, nil
}

func (s *Service) SaveUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "services.pcClub.user.SaveUser"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors2.WithMessage(err, op, "failed hash password")
	}

	id, err := s.userOwner.SaveUser(
		ctx,
		&models.User{
			Email:    email,
			Password: passHash,
		})
	if err != nil {
		return 0, errors2.WithMessage(HandleStorageError(err), op, "failed to save user in mssql")
	}

	return id, nil
}

func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "services.pcClub.user.Email"

	user, err := s.userProvider.UserByEmail(ctx, email)
	if errors.Is(err, mssql.ErrNotFound) {
		return 0, errors2.WithMessage(ErrInvalidCredentials, op, "failed to get user from mssql")
	}
	if err != nil {
		return 0, errors2.WithMessage(HandleStorageError(err), op, "failed to get user from mssql")
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return 0, errors2.WithMessage(HandleStorageError(err), op, "failed to compare password")
	}

	return user.UserID, nil
}

func (s *Service) DeleteUser(
	ctx context.Context,
	uid int64,
) error {
	const op = "services.pcClub.user.DeleteUser"

	err := s.userOwner.DeleteUser(ctx, uid)
	if err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to delete user from mssql")
	}

	return nil
}

func (s *Service) IsAdmin(
	ctx context.Context,
	uid int64,
) error {
	const op = "services.pcClub.user.IsAdmin"

	role, err := s.userProvider.UserRole(ctx, uid)
	if err != nil {
		return errors2.WithMessage(HandleStorageError(err), op, "failed to get user role from mssql")
	}

	if role != s.cfg.AdminRoleName {
		return errors2.WithMessage(ErrAccessDenied, op, "user is not admin")
	}

	return nil
}
