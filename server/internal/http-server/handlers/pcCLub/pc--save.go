package pcCLub

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"server/internal/lib/api/response"
	"server/internal/lib/logger/sl"
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

type SavePcRequest struct {
	TypeId int64 `json:"type_id" validate:"required,numeric"`
	RoomId int64 `json:"room_id" validate:"required,numeric"`
	Row    int   `json:"row" validate:"required,numeric"`
	Place  int   `json:"place" validate:"required,numeric"`
}

func (a *API) SavePcType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc.SavePcType"

		log := a.Log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

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
			if errors.Is(err, pc.ErrPcTypeAlreadyExists) {
				log.Warn("pc type already exists", sl.Err(err))
				response.AlreadyExists(w, "pc type already exists")
				return
			}

			log.Error("failed to save pc type", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, "saved success")
	}
}

func (a *API) SavePc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc.SavePc"

		log := a.Log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

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
		)
		if errors.Is(err, pc.ErrPcAlreadyExists) {
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
