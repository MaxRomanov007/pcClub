package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/api/logger/sl"
	"server/internal/lib/api/request"
	"server/internal/lib/api/response"
	"server/internal/models"
	"server/internal/services/pcClub/components"
)

type RamsRequest struct {
	TypeId int64 `get:"type-id" validate:"required,min=1"`
}

type SaveRamTypeRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type SaveRamRequest struct {
	TypeId   int64 `json:"type_id" validate:"required,min=1"`
	Capacity int   `json:"capacity" validate:"required,numeric,min=0"`
}

type DeleteRamTypeRequest struct {
	TypeId int64 `json:"type_id" validate:"required,min=1"`
}

type DeleteRamRequest struct {
	RamId int64 `json:"ram_id" validate:"required,min=1"`
}

func (a *API) RamTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.ram.RamTypes"

		log := a.log(op, r)

		types, err := a.ComponentsService.Ram.RamTypes(r.Context())
		if err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to get ram types", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, types)
	}
}

func (a *API) Rams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.ram.Rams"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateGETRequest[RamsRequest](w, r, log)
		if !ok {
			return
		}

		rams, err := a.ComponentsService.Ram.Rams(r.Context(), req.TypeId)
		if err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to get rams", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, rams)
	}
}

func (a *API) SaveRamType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.ram.SaveRamType"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[SaveRamTypeRequest](w, r, log)
		if !ok {
			return
		}

		ramType := models.RAMType{
			Name: req.Name,
		}
		if _, err := a.ComponentsService.Ram.SaveRamType(r.Context(), &ramType); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to save ram type", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, ramType)
	}
}

func (a *API) SaveRam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.ram.SaveRam"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[SaveRamRequest](w, r, log)
		if !ok {
			return
		}

		ram := models.RAM{
			RAMTypeID: req.TypeId,
			Capacity:  req.Capacity,
		}
		if _, err := a.ComponentsService.Ram.SaveRam(r.Context(), &ram); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to save ram", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, ram)
	}
}

func (a *API) DeleteRamType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.ram.DeleteRamType"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[DeleteRamTypeRequest](w, r, log)
		if !ok {
			return
		}

		if err := a.ComponentsService.Ram.DeleteRamType(r.Context(), req.TypeId); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to delete ram type", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) DeleteRam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.ram.DeleteRam"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[DeleteRamRequest](w, r, log)
		if !ok {
			return
		}

		if err := a.ComponentsService.Ram.DeleteRam(r.Context(), req.RamId); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to delete ram", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}
