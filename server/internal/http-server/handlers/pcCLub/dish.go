package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/api/logger/sl"
	"server/internal/lib/api/request"
	"server/internal/lib/api/response"
	"server/internal/models"
	dishService "server/internal/services/pcClub/dish"
)

type DishesRequest struct {
	Limit  int `get:"limit" validate:"omitempty,min=0"`
	Offset int `get:"offset" validate:"omitempty,min=0"`
}

type DishRequest struct {
	DishId int64 `get:"dish-id,true" validate:"required,min=1"`
}

type SaveDishRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Calories    int16   `json:"calories" validate:"required,min=0"`
	Cost        float32 `json:"cost" validate:"required,min=0"`
	Description string  `json:"description" validate:"omitempty"`
}

type UpdateDishRequest struct {
	DishId      int64   `json:"dish_id" validate:"required,min=1"`
	StatusId    int64   `json:"status_id" validate:"omitempty,min=1"`
	Name        string  `json:"name" validate:"omitempty,max=255"`
	Calories    int16   `json:"calories" validate:"omitempty,min=0"`
	Cost        float32 `json:"cost" validate:"omitempty,min=0"`
	Description string  `json:"description" validate:"omitempty"`
}

type DeleteDishRequest struct {
	DishId int64 `json:"dish_id" validate:"required,min=1"`
}

func (a *API) Dishes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.dish.Dishes"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateGETRequest[DishesRequest](w, r, log)
		if !ok {
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

		req, ok := request.DecodeAndValidateGETRequest[DishRequest](w, r, log)
		if !ok {
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

		req, ok := request.DecodeAndValidateJSONRequest[SaveDishRequest](w, r, log)
		if !ok {
			return
		}

		dish := models.Dish{
			Name:        req.Name,
			Calories:    req.Calories,
			Cost:        req.Cost,
			Description: req.Description,
		}
		if _, err := a.DishService.SaveDish(r.Context(), &dish); err != nil {
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

		render.JSON(w, r, dish)
	}
}

func (a *API) UpdateDish() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.dish.UpdateDish"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[UpdateDishRequest](w, r, log)
		if !ok {
			return
		}

		dish := models.Dish{
			DishStatusID: req.StatusId,
			Name:         req.Name,
			Calories:     req.Calories,
			Cost:         req.Cost,
			Description:  req.Description,
		}
		if err := a.DishService.UpdateDish(r.Context(), req.DishId, &dish); err != nil {
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

		render.JSON(w, r, dish)
	}
}

func (a *API) DeleteDish() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.dish.DeleteDish"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[DeleteDishRequest](w, r, log)
		if !ok {
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
