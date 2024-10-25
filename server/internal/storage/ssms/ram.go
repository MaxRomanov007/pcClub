package ssms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

func (s *Storage) RamTypes(
	ctx context.Context,
) ([]models.RamType, error) {
	const op = "storage.ssms.ram.RamTypes"

	query, args, err := squirrel.Select("*").From("ram_types").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var types []models.RamType
	if err := s.db.SelectContext(ctx, &types, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to select ram types: %w", op, handleError(err))
	}

	return types, nil
}

func (s *Storage) Rams(
	ctx context.Context,
	typeId int64,
) ([]models.Ram, error) {
	const op = "storage.ssms.ram.Rams"

	query, args, err := squirrel.
		Select("*").
		From("ram").
		Where(squirrel.Eq{
			"ram_type_id": typeId,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var rams []models.Ram
	if err := s.db.SelectContext(ctx, &rams, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to select rams: %w", op, handleError(err))
	}

	return rams, nil
}

func (s *Storage) SaveRamType(
	ctx context.Context,
	name string,
) error {
	const op = "storage.ssms.ram.SaveRamType"

	query, args, err := squirrel.Insert("ram_types").Values(name).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to save ram type: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) SaveRam(
	ctx context.Context,
	ram models.Ram,
) error {
	const op = "storage.ssms.ram.SaveRam"

	query, args, err := squirrel.
		Insert("ram").
		Columns(
			"ram_type_id",
			"model",
		).
		Values(
			ram.RamTypeID,
			ram.Capacity,
		).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to save ram: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeleteRamType(
	ctx context.Context,
	typeId int64,
) error {
	const op = "storage.ssms.ram.DeleteRamType"

	query, args, err := squirrel.
		Delete("ram_types").
		Where(squirrel.Eq{
			"ram_type_id": typeId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to delete ram type: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeleteRam(
	ctx context.Context,
	ramId int64,
) error {
	const op = "storage.ssms.ram.DeleteRam"

	query, args, err := squirrel.
		Delete("ram").
		Where(squirrel.Eq{
			"ram_id": ramId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to delete ram: %w", op, handleError(err))
	}

	return nil
}
