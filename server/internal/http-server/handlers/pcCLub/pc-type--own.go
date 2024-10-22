package pcCLub

import (
	"errors"
	"net/http"
	"server/internal/lib/logger/sl"
	"server/internal/lib/response"
	"server/internal/models"
	"server/internal/services/pcClub/pc"
)

type SavePcTypeRequest struct {
	Name        string                `json:"name" validate:"required"`
	Description string                `json:"description" validate:"omitempty,max=255"`
	Processor   *models.ProcessorData `json:"processor" validate:"required,dive"`
	VideoCard   *models.VideoCardData `json:"video_card" validate:"required,dive"`
	Monitor     *models.MonitorData   `json:"monitor" validate:"required,dive"`
	Ram         *models.RamData       `json:"ram" validate:"required,dive"`
}

type UpdatePcTypeRequest struct {
	TypeId      int64                 `json:"id" validate:"required,numeric"`
	Name        string                `json:"name" validate:"required"`
	Description string                `json:"description" validate:"omitempty,max=255"`
	Processor   *models.ProcessorData `json:"processor" validate:"required,dive"`
	VideoCard   *models.VideoCardData `json:"video_card" validate:"required,dive"`
	Monitor     *models.MonitorData   `json:"monitor" validate:"required,dive"`
	Ram         *models.RamData       `json:"ram" validate:"required,dive"`
}

type DeletePcTypeRequest struct {
	PcTypeId int64 `json:"pc_type_id" validate:"required,numeric"`
}

func (a *API) SavePcType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pcType.SavePcType"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SavePcTypeRequest
		if !a.decodeAndValidateRequest(w, r, log, &req) {
			return
		}

		if err := a.PcService.SavePcType(
			r.Context(),
			req.Name,
			req.Description,
			req.Processor,
			req.VideoCard,
			req.Monitor,
			req.Ram,
		); err != nil {
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
	}
}

func (a *API) UpdatePcType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pcType.UpdatePcType"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req UpdatePcTypeRequest
		if !a.decodeAndValidateRequest(w, r, log, &req) {
			return
		}

		if err := a.PcService.UpdatePcType(
			r.Context(),
			req.TypeId,
			req.Name,
			req.Description,
			req.Processor,
			req.VideoCard,
			req.Monitor,
			req.Ram,
		); err != nil {
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
	}
}

func (a *API) DeletePcType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pcType.DeletePcType"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req DeletePcTypeRequest
		if !a.decodeAndValidateRequest(w, r, log, &req) {
			return
		}

		err := a.PcService.DeletePcType(r.Context(), req.PcTypeId)
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
