package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"net/http"
	"server/internal/lib/api/response"
	"server/internal/lib/logger/sl"
	"server/internal/services/pcClub/pc"
	"strconv"
)

type PcsRequest struct {
	typeId      string `validate:"require,number"`
	isAvailable string `validate:"omitempty,boolean"`
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

func (a *API) Pcs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc"

		log := a.log(op, r)

		req := &PcsRequest{
			typeId:      r.URL.Query().Get("type-id"),
			isAvailable: r.URL.Query().Get("is-available"),
		}
		if err := validator.New().Struct(req); err != nil {
			var validErr validator.ValidationErrors
			if ok := errors.Is(err, &validErr); ok {
				log.Warn("invalid request", sl.Err(err))
				response.ValidationFailed(w, validErr)
				return
			}

			log.Error("validation failed", sl.Err(err))
			response.Internal(w)
			return
		}

		if req.typeId == "" {
			req.typeId = "0"
		}
		if req.isAvailable == "" {
			req.isAvailable = "0"
		}

		typeId, err := strconv.ParseInt(req.typeId, 10, 64)
		if err != nil {
			log.Error("failed to parse type id", sl.Err(err))
			response.Internal(w)
			return
		}

		isFree, err := strconv.ParseBool(req.isAvailable)
		if err != nil {
			log.Error("failed to parse is free flag", sl.Err(err))
			response.Internal(w)
			return
		}

		pcs, err := a.PcService.Pcs(r.Context(), typeId, isFree)
		if errors.Is(err, pc.ErrNotFound) {
			log.Warn("pc type not found", sl.Err(err))
			response.NotFound(w, "pc type not found")
			return
		}
		if err != nil {
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
		if !a.decodeAndValidateRequest(w, r, log, &req) {
			return
		}

		err := a.PcService.SavePc(
			r.Context(),
			req.TypeId,
			req.RoomId,
			req.Row,
			req.Place,
			req.Description,
		)
		if errors.Is(err, pc.ErrAlreadyExists) {
			log.Warn("pc already exists", sl.Err(err))
			response.AlreadyExists(w, "pc already exists")
			return
		}
		if err != nil {
			log.Error("failed to save pc")
			response.Internal(w)
			return
		}

		render.JSON(w, r, "saved success")
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
		if !a.decodeAndValidateRequest(w, r, log, &req) {
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
			//TODO: make pc type error and response handle method
			if errors.Is(err, pc.ErrConstraint) {
				log.Warn("pc row or place greater than room", sl.Err(err))
				response.AlreadyExists(w, "pc row or place greater than room")
				return
			}
			if errors.Is(err, pc.ErrAlreadyExists) {
				log.Warn("pc already exists", sl.Err(err))
				response.AlreadyExists(w, "pc already exists")
				return
			}
			if errors.Is(err, pc.ErrReferenceNotExists) {
				log.Warn("pc already exists", sl.Err(err))
				response.AlreadyExists(w, "reference data doesnt exists")
				return
			}
			if errors.Is(err, pc.ErrNotFound) {
				log.Warn("pc not found", sl.Err(err))
				response.NotFound(w, "pc not found")
				return
			}

			log.Error("failed to update pc", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}
