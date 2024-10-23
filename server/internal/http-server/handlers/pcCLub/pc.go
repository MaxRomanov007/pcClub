package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/logger/sl"
	"server/internal/lib/response"
	"server/internal/services/pcClub/pc"
)

type PcsRequest struct {
	TypeId      int64 `validate:"require,number,min=1" get:"type-id"`
	IsAvailable bool  `validate:"omitempty,boolean" get:"is-available"`
}

type SavePcRequest struct {
	TypeId      int64  `json:"type_id" validate:"required,numeric"`
	RoomId      int64  `json:"room_id" validate:"required,numeric"`
	Row         int    `json:"row" validate:"required,numeric"`
	Place       int    `json:"place" validate:"required,numeric"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdatePcRequest struct {
	PcId        int64  `json:"pc_id" validate:"required,numeric"`
	TypeId      int64  `json:"type_id" validate:"omitempty,numeric"`
	RoomId      int64  `json:"room_id" validate:"omitempty,numeric"`
	StatusId    int64  `json:"status_id" validate:"omitempty,numeric"`
	Row         int    `json:"row" validate:"omitempty,numeric"`
	Place       int    `json:"place" validate:"omitempty,numeric"`
	Description string `json:"description" validate:"omitempty"`
}

type DeletePcRequest struct {
	PcId int64 `json:"pc_id" validate:"required,numeric"`
}

func (a *API) Pcs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc"

		log := a.log(op, r)

		var req PcsRequest
		if !a.decodeAndValidateGETRequest(w, r, log, &req) {
			return
		}

		pcs, err := a.PcService.Pcs(r.Context(), req.TypeId, req.IsAvailable)
		if err != nil {
			var pcErr *pc.Error
			if ok := errors.As(err, &pcErr); ok {
				log.Warn("pc error", sl.Err(err))
				response.PcError(w, pcErr)
				return
			}
			log.Error("failed to get pcs", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, pcs)
	}
}

func (a *API) SavePc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc.SavePc"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SavePcRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.PcService.SavePc(
			r.Context(),
			req.TypeId,
			req.RoomId,
			req.Row,
			req.Place,
			req.Description,
		); err != nil {
			var pcErr *pc.Error
			if ok := errors.As(err, &pcErr); ok {
				log.Warn("pc error", sl.Err(err))
				response.PcError(w, pcErr)
				return
			}
			log.Error("failed to save pc", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) UpdatePc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc.SavePc"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req UpdatePcRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.PcService.UpdatePc(
			r.Context(),
			req.PcId,
			req.TypeId,
			req.RoomId,
			req.StatusId,
			req.Row,
			req.Place,
			req.Description,
		); err != nil {
			var pcErr *pc.Error
			if ok := errors.As(err, &pcErr); ok {
				log.Warn("pc error", sl.Err(err))
				response.PcError(w, pcErr)
				return
			}

			log.Error("failed to update pc", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) DeletePc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc.SavePc"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req DeletePcRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.PcService.DeletePc(r.Context(), req.PcId); err != nil {
			var pcErr *pc.Error
			if ok := errors.As(err, &pcErr); ok {
				log.Warn("pc error", sl.Err(err))
				response.PcError(w, pcErr)
				return
			}
			log.Error("failed to delete pc", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}
