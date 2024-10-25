package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/logger/sl"
	"server/internal/lib/response"
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

		var req RamsRequest
		if !a.decodeAndValidateGETRequest(w, r, log, &req) {
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

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SaveRamTypeRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.Ram.SaveRamType(r.Context(), req.Name); err != nil {
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
	}
}

func (a *API) SaveRam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.ram.SaveRam"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SaveRamRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.Ram.SaveRam(r.Context(), models.Ram{
			RamTypeID: req.TypeId,
			Capacity:  req.Capacity,
		}); err != nil {
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
	}
}

func (a *API) DeleteRamType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.ram.DeleteRamType"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req DeleteRamTypeRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
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

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req DeleteRamRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
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
