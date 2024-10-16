package sqlServer

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

func (s *Storage) User(
	ctx context.Context,
	uid int64,
) (models.UserData, error) {
	const op = "storage.sqlServer.user.User"

	stmt, args, err := squirrel.Select("user_id", "email", "balance").
		From("users").
		Where(squirrel.Eq{
			"user_id": uid,
		}).ToSql()
	if err != nil {
		return models.UserData{}, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	var user models.UserData
	err = s.db.GetContext(ctx, &user, stmt, args...)
	if err != nil {
		return models.UserData{}, fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return user, nil
}

func (s *Storage) UserByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "storage.sqlServer.user.UserByEmail"

	stmt, args, err := squirrel.Select("*").
		From("users").
		Where(squirrel.Eq{
			"email": email,
		}).ToSql()
	if err != nil {
		return models.User{}, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	var user models.User
	err = s.db.GetContext(ctx, &user, stmt, args...)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return user, nil
}

func (s *Storage) SaveUser(
	ctx context.Context,
	user models.User,
) (int64, error) {
	const op = "storage.sqlServer.user.SaveUser"

	stmt, args, err := squirrel.Insert("users").
		Columns("email", "password").
		Values(user.Email, user.Password).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	_, err = s.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	stmt, args, err = squirrel.Select("user_id").
		From("users").
		Where(squirrel.Eq{
			"email": user.Email,
		}).
		ToSql()

	stmt = replacePositionalParams(stmt, args)

	var id int64
	err = s.db.GetContext(ctx, &id, stmt, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get id: %w", op, handleError(err))
	}

	return id, nil
}

func (s *Storage) DeleteUser(
	ctx context.Context,
	uid int64,
) error {
	const op = "storage.sqlServer.user.DeleteUser"

	stmt, args, err := squirrel.
		Delete("users").
		Where(squirrel.Eq{
			"user_id": uid,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	_, err = s.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) RefreshVersion(
	ctx context.Context,
	uid int64,
) (int64, error) {
	const op = "storage.sqlServer.user.RefreshVersion"

	stmt, args, err := squirrel.Select("refresh_token_version").
		From("users").
		Where(squirrel.Eq{
			"user_id": uid,
		}).ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	var user models.User
	err = s.db.GetContext(ctx, &user, stmt, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return user.RefreshTokenVersion, nil
}

func (s *Storage) UpdateRefreshVersion(
	ctx context.Context,
	uid int64,
	version int64,
) error {
	const op = "storage.sqlServer.user.UpdateRefreshVersion"

	stmt, args, err := squirrel.Update("users").
		Set("refresh_token_version", version).
		Where(squirrel.Eq{
			"user_id": uid,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	_, err = s.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) UserRole(
	ctx context.Context,
	uid int64,
) (string, error) {
	const op = "storage.sqlServer.user.UserRole"

	stmt, args, err := squirrel.
		Select("user_roles.name").
		From("user_roles").
		Join("users ON users.user_role_id = user_roles.user_role_id").
		Where(squirrel.Eq{
			"users.user_id": uid,
		}).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	var role string
	if err := s.db.GetContext(ctx, &role, stmt, args...); err != nil {
		return "", fmt.Errorf("%s: failed to get role: %w", op, handleError(err))
	}

	return role, nil
}
