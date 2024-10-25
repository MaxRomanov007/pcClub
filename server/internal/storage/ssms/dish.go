package ssms

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

type DishDataNullable struct {
	Id          int64            `db:"dish_id" json:"id"`
	StatusId    int64            `db:"dish_status_id" json:"status_id"`
	Status      string           `db:"status" json:"status"`
	Name        string           `db:"name" json:"name"`
	Calories    int16            `db:"calories" json:"calories"`
	Cost        float64          `db:"cost" json:"cost"`
	Description sql.Null[string] `db:"description" json:"description"`
}

func (d DishDataNullable) Parse() models.DishData {
	return models.DishData{
		Id:          d.Id,
		StatusId:    d.StatusId,
		Status:      d.Status,
		Name:        d.Name,
		Calories:    d.Calories,
		Cost:        d.Cost,
		Description: d.Description.V,
	}
}

func (s *Storage) Dishes(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]models.DishData, error) {
	const op = "storage.ssms.dish.Dishes"

	builder := squirrel.Select(
		"dishes.dish_id",
		"dishes.name",
		"dishes.calories",
		"dishes.cost",
		"dishes.description",
		"dish_statuses.name AS status",
	).
		From("dishes").
		Join("dish_statuses ON dish_statuses.dish_status_id = dishes.dish_status_id")

	if limit > 0 || offset > 0 {
		builder = builder.
			Column("ROW_NUMBER() OVER (ORDER BY dishes.dish_id) AS nom").
			Prefix("WITH num_row AS (").
			Suffix(")")

		query, args, err := builder.ToSql()
		if err != nil {
			return nil, fmt.Errorf("%s: failed to generate inner query: %w", op, err)
		}

		builder = squirrel.Select(
			"dish_id",
			"name",
			"calories",
			"cost",
			"description",
			"status",
		).
			From("num_row").
			Where("nom BETWEEN ? AND ?", offset+1, limit+offset).
			Prefix(query, args...)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to generate query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var dishesNull []DishDataNullable
	if err := s.db.SelectContext(ctx, &dishesNull, query, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to select dishes: %w", op, handleError(err))
	}

	dishes := make([]models.DishData, len(dishesNull))
	for i, dish := range dishesNull {
		dishes[i] = dish.Parse()
	}

	return dishes, nil
}

func (s *Storage) Dish(
	ctx context.Context,
	dishId int64,
) (models.DishData, error) {
	const op = "storage.ssms.dish.Dish"

	query, args, err := squirrel.Select(
		"dishes.dish_id",
		"dishes.name",
		"dishes.calories",
		"dishes.cost",
		"dishes.description",
		"dish_statuses.name AS status",
	).
		From("dishes").
		Where(squirrel.Eq{
			"dishes.dish_id": dishId,
		}).
		Join("dish_statuses ON dish_statuses.dish_status_id = dishes.dish_status_id").
		ToSql()
	if err != nil {
		return models.DishData{}, fmt.Errorf("%s: failed to generate query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	var dish DishDataNullable
	if err := s.db.GetContext(ctx, &dish, query, args...); err != nil {
		return models.DishData{}, fmt.Errorf("%s: failed to get dish: %w", op, handleError(err))
	}

	return dish.Parse(), nil
}

func (s *Storage) SaveDish(
	ctx context.Context,
	dish models.DishData,
) error {
	const op = "storage.ssms.dish.SaveDish"

	query, args, err := squirrel.Insert("dishes").Columns(
		"name",
		"calories",
		"cost",
		"description",
	).
		Values(
			dish.Name,
			dish.Calories,
			dish.Cost,
			nullString(dish.Description),
		).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to generate query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) UpdateDish(
	ctx context.Context,
	dish models.DishData,
) error {
	const op = "storage.ssms.dish.UpdateDish"

	builder := squirrel.Update("dishes").Where(squirrel.Eq{
		"dish_id": dish.Id,
	})
	builder = setIfNotZero(builder, "name", dish.Name)
	builder = setIfNotZero(builder, "dish_status_id", dish.StatusId)
	builder = setIfNotZero(builder, "calories", dish.Calories)
	builder = setIfNotZero(builder, "cost", dish.Cost)
	builder = setIfNotZeroNullString(builder, "description", dish.Description)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to generate query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) DeleteDish(
	ctx context.Context,
	dishId int64,
) error {
	const op = "storage.ssms.dish.DeleteDish"

	query, args, err := squirrel.
		Delete("dishes").
		Where(squirrel.Eq{
			"dish_id": dishId,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to generate query: %w", op, err)
	}

	query = replacePositionalParams(query, args)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", op, handleError(err))
	}

	return nil
}
