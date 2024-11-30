package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/api/logger/sl"
	"server/internal/lib/api/request"
	"server/internal/lib/api/response"
	"server/internal/models"
	"server/internal/services/pcClub/pc"
)

type PcTypesRequest struct {
	Limit  int `validate:"omitempty,number,min=0" get:"limit"`
	Offset int `validate:"omitempty,number,min=0" get:"offset"`
}

type PcTypeRequest struct {
	TypeId int64 `validate:"required,number,min=1" get:"type-id, true"`
}

type SavePcTypeRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"omitempty,max=255"`
	HourCost    float32 `json:"hour_cost" validate:"required,min=1"`
	ProcessorID int64   `json:"processor_id" validate:"required,min=1"`
	VideoCardID int64   `json:"video_card_id" validate:"required,min=1"`
	MonitorID   int64   `json:"monitor_id" validate:"required,min=1"`
	RamID       int64   `json:"ram_id" validate:"required,min=1"`
}

type UpdatePcTypeRequest struct {
	TypeID      int64   `json:"id" validate:"required,numeric"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"omitempty,max=255"`
	HourCost    float32 `json:"hour_cost" validate:"required,min=1"`
	ProcessorID int64   `json:"processor_id" validate:"required,min=1"`
	VideoCardID int64   `json:"video_card_id" validate:"required,min=1"`
	MonitorID   int64   `json:"monitor_id" validate:"required,min=1"`
	RamID       int64   `json:"ram_id" validate:"required,min=1"`
}

type DeletePcTypeRequest struct {
	PcTypeId int64 `json:"pc_type_id" validate:"required,numeric"`
}

func (a *API) PcTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pcType.PcTypes"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateGETRequest[PcTypesRequest](w, r, log)
		if !ok {
			return
		}

		pcTypes, err := a.PcTypeService.PcTypes(r.Context(), req.Limit, req.Offset)
		if err != nil {
			var pcErr *pc.Error
			if ok := errors.As(err, &pcErr); ok {
				log.Warn("pc error", sl.Err(err))
				response.PcError(w, pcErr)
				return
			}
			log.Error("failed to get pc types", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, pcTypes)
	}
}

func (a *API) PcType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pcType.PcType"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateGETRequest[PcTypeRequest](w, r, log)
		if !ok {
			return
		}

		pcType, err := a.PcTypeService.PcType(r.Context(), req.TypeId)
		if err != nil {
			var pcErr *pc.Error
			if ok := errors.As(err, &pcErr); ok {
				log.Warn("pc error", sl.Err(err))
				response.PcError(w, pcErr)
				return
			}
			log.Error("failed to get pc type", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, pcType)
	}
}

func (a *API) SavePcType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pcType.SavePcType"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[SavePcTypeRequest](w, r, log)
		if !ok {
			return
		}

		pcType := models.PcType{
			Name:        req.Name,
			Description: req.Description,
			HourCost:    req.HourCost,
			ProcessorID: req.ProcessorID,
			VideoCardID: req.VideoCardID,
			MonitorID:   req.MonitorID,
			RAMID:       req.RamID,
		}
		if _, err := a.PcTypeService.SavePcType(r.Context(), &pcType); err != nil {
			var pcErr *pc.Error
			if ok := errors.As(err, &pcErr); ok {
				log.Warn("pc error", sl.Err(err))
				response.PcError(w, pcErr)
				return
			}
			log.Error("failed to save pc type", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, pcType)
	}
}

func (a *API) UpdatePcType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pcType.UpdatePcType"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[UpdatePcTypeRequest](w, r, log)
		if !ok {
			return
		}

		pcType := models.PcType{
			Name:        req.Name,
			Description: req.Description,
			HourCost:    req.HourCost,
			ProcessorID: req.ProcessorID,
			VideoCardID: req.VideoCardID,
			MonitorID:   req.MonitorID,
			RAMID:       req.RamID,
		}
		if err := a.PcTypeService.UpdatePcType(r.Context(), req.TypeID, &pcType); err != nil {
			var pcErr *pc.Error
			if ok := errors.As(err, &pcErr); ok {
				log.Warn("pc error", sl.Err(err))
				response.PcError(w, pcErr)
				return
			}
			log.Error("failed to save pc type", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, pcType)
	}
}

func (a *API) DeletePcType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pcType.DeletePcType"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[DeletePcTypeRequest](w, r, log)
		if !ok {
			return
		}

		err := a.PcTypeService.DeletePcType(r.Context(), req.PcTypeId)
		if err != nil {
			var pcErr *pc.Error
			if ok := errors.As(err, &pcErr); ok {
				log.Warn("pc error", sl.Err(err))
				response.PcError(w, pcErr)
				return
			}
			log.Error("failed to delete pc type", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}
