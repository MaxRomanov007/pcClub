package user

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"server/internal/models"
	"server/internal/storage/ssms"
)

func (s *Service) User(
	ctx context.Context,
	uid int64,
) (models.UserData, error) {
	const op = "services.pcClub.user.User"

	user, err := s.userProvider.User(ctx, uid)
	if errors.Is(err, ssms.ErrNotFound) {
		return models.UserData{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	if err != nil {
		return models.UserData{}, fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	return user, nil
}

func (s *Service) UserByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "services.pcClub.user.UserByEmail"

	user, err := s.userProvider.UserByEmail(ctx, email)
	if errors.Is(err, ssms.ErrNotFound) {
		return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	if err != nil {
		return models.User{}, fmt.Errorf("%s: failed to get user: %w", op, err)
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
		return 0, fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	id, err := s.userOwner.SaveUser(
		ctx,
		models.User{
			Email:    email,
			Password: passHash,
		})
	if errors.Is(err, ssms.ErrAlreadyExists) {
		return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: failed to save user: %w", op, err)
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
	if errors.Is(err, ssms.ErrNotFound) {
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: failed to compare password: %w", op, err)
	}

	return user.UserID, nil
}

func (s *Service) DeleteUser(
	ctx context.Context,
	uid int64,
) error {
	const op = "services.pcClub.user.DeleteUser"

	err := s.userOwner.DeleteUser(ctx, uid)
	if errors.Is(err, ssms.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: failed to delete user: %w", op, err)
	}

	return nil
}

func (s *Service) IsAdmin(
	ctx context.Context,
	uid int64,
) error {
	const op = "services.pcClub.user.IsAdmin"

	role, err := s.userProvider.UserRole(ctx, uid)
	if errors.Is(err, ssms.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: failed to get role: %w", op, err)
	}

	if role != s.cfg.AdminRoleName {
		return fmt.Errorf("%s: %w", op, ErrAccessDenied)
	}

	return nil
}
