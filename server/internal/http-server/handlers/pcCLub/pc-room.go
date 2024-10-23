package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/logger/sl"
	"server/internal/lib/response"
	"server/internal/models"
	pcRoom2 "server/internal/services/pcClub/pcRoom"
)

type PcRoomRequest struct {
	RoomId int64 `get:"room-id,true" validate:"required,numeric,min=1"`
}

type SavePcRoomRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Rows        int    `json:"rows" validate:"required,numeric,min=1"`
	Places      int    `json:"places" validate:"required,numeric,min=1"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdatePcRoomRequest struct {
	RoomId      int64  `json:"room_id" validate:"required,numeric,min=1"`
	Name        string `json:"name" validate:"omitempty,max=255"`
	Rows        int    `json:"rows" validate:"omitempty,numeric,min=1"`
	Places      int    `json:"places" validate:"omitempty,numeric,min=1"`
	Description string `json:"description" validate:"omitempty"`
}

type DeletePcRoomRequest struct {
	RoomId int64 `json:"room_id" validate:"required,numeric,min=1"`
}

func (a *API) PcRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.PcRoom"

		log := a.log(op, r)

		var req PcRoomRequest
		if !a.decodeAndValidateGETRequest(w, r, log, &req) {
			return
		}

		pcRoom, err := a.PcRoomService.PcRoom(r.Context(), req.RoomId)
		if err != nil {
			var serviceErr *pcRoom2.Error
			if errors.As(err, &serviceErr) {
				log.Warn("pc room error", sl.Err(err))
				response.PcRoomError(w, serviceErr)
				return
			}
			log.Error("failed to get pc room", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, pcRoom)
	}
}

func (a *API) SavePcRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.savePcRoom"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SavePcRoomRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.PcRoomService.SavePcRoom(
			r.Context(),
			models.PcRoom{
				Name:        req.Name,
				Rows:        req.Rows,
				Places:      req.Places,
				Description: req.Description,
			},
		); err != nil {
			var serviceErr *pcRoom2.Error
			if errors.As(err, &serviceErr) {
				log.Warn("pc room error", sl.Err(err))
				response.PcRoomError(w, serviceErr)
				return
			}
			log.Error("failed to save pc room", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) UpdatePcRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.UpdatePcRoom"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req UpdatePcRoomRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.PcRoomService.UpdatePcRoom(
			r.Context(),
			models.PcRoom{
				PcRoomID:    req.RoomId,
				Name:        req.Name,
				Rows:        req.Rows,
				Places:      req.Places,
				Description: req.Description,
			},
		); err != nil {
			var serviceErr *pcRoom2.Error
			if errors.As(err, &serviceErr) {
				log.Warn("pc room error", sl.Err(err))
				response.PcRoomError(w, serviceErr)
				return
			}
			log.Error("failed to update pc room", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) DeletePcRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.DeletePcRoom"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req DeletePcRoomRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.PcRoomService.DeletePcRoom(r.Context(), req.RoomId); err != nil {
			var serviceErr *pcRoom2.Error
			if errors.As(err, &serviceErr) {
				log.Warn("pc room error", sl.Err(err))
				response.PcRoomError(w, serviceErr)
				return
			}
			log.Error("failed to update pc room", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}
