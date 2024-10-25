package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/logger/sl"
	"server/internal/lib/response"
	"server/internal/models"
	dishService "server/internal/services/pcClub/dish"
)

type DishesRequest struct {
	Limit  int64 `get:"limit" validate:"omitempty,min=0"`
	Offset int64 `get:"offset" validate:"omitempty,min=0"`
}

type DishRequest struct {
	DishId int64 `get:"dish-id,true" validate:"required,min=1"`
}

type SaveDishRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Calories    int16   `json:"calories" validate:"required,min=0"`
	Cost        float64 `json:"cost" validate:"required,min=0"`
	Description string  `json:"description" validate:"omitempty"`
}

type UpdateDishRequest struct {
	DishId      int64   `json:"dish_id" validate:"required,min=1"`
	StatusId    int64   `json:"status_id" validate:"omitempty,min=1"`
	Name        string  `json:"name" validate:"omitempty,max=255"`
	Calories    int16   `json:"calories" validate:"omitempty,min=0"`
	Cost        float64 `json:"cost" validate:"omitempty,min=0"`
	Description string  `json:"description" validate:"omitempty"`
}

type DeleteDishRequest struct {
	DishId int64 `json:"dish_id" validate:"required,min=1"`
}

func (a *API) Dishes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.dish.Dishes"

		log := a.log(op, r)

		var req DishesRequest
		if !a.decodeAndValidateGETRequest(w, r, log, &req) {
			return
		}

		dishes, err := a.DishService.Dishes(r.Context(), req.Limit, req.Offset)
		if err != nil {
			var serviceErr *dishService.Error
			if errors.As(err, &serviceErr) {
				log.Warn("dish error", sl.Err(err))
				response.DishError(w, serviceErr)
				return
			}
			log.Error("failed to get dishes", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, dishes)
	}
}

func (a *API) Dish() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.dish.Dish"

		log := a.log(op, r)

		var req DishRequest
		if !a.decodeAndValidateGETRequest(w, r, log, &req) {
			return
		}

		dish, err := a.DishService.Dish(r.Context(), req.DishId)
		if err != nil {
			var serviceErr *dishService.Error
			if errors.As(err, &serviceErr) {
				log.Warn("dish error", sl.Err(serviceErr))
				response.DishError(w, serviceErr)
				return
			}
			log.Error("failed to get dish", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, dish)
	}
}

func (a *API) SaveDish() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.dish.SaveDish"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SaveDishRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.DishService.SaveDish(r.Context(), models.DishData{
			Name:        req.Name,
			Calories:    req.Calories,
			Cost:        req.Cost,
			Description: req.Description,
		}); err != nil {
			var serviceErr *dishService.Error
			if errors.As(err, &serviceErr) {
				log.Warn("dish error", sl.Err(serviceErr))
				response.DishError(w, serviceErr)
				return
			}
			log.Error("failed to save dish", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) UpdateDish() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.dish.UpdateDish"

		log := a.log(op, r)

		var req UpdateDishRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.DishService.UpdateDish(r.Context(), models.DishData{
			Id:          req.DishId,
			StatusId:    req.StatusId,
			Name:        req.Name,
			Calories:    req.Calories,
			Cost:        req.Cost,
			Description: req.Description,
		}); err != nil {
			var serviceErr *dishService.Error
			if errors.As(err, &serviceErr) {
				log.Warn("dish error", sl.Err(err))
				response.DishError(w, serviceErr)
				return
			}
			log.Error("failed to update dish", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) DeleteDish() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.dish.DeleteDish"

		log := a.log(op, r)

		var req DeleteDishRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.DishService.DeleteDish(r.Context(), req.DishId); err != nil {
			var serviceErr *dishService.Error
			if errors.As(err, &serviceErr) {
				log.Warn("dish error", sl.Err(err))
				response.DishError(w, serviceErr)
				return
			}
			log.Error("failed to delete dish", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}
