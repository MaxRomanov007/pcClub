package pcCLub

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
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

type PcTypesRequest struct {
	limit  string `validate:"omitempty,number"`
	offset string `validate:"omitempty,number"`
}

type PcTypeRequest struct {
	typeId string `validate:"require,number"`
}

func (a *API) Pcs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc"

		log := a.Log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req := &PcsRequest{
			typeId:      r.URL.Query().Get("type_id"),
			isAvailable: r.URL.Query().Get("is_available"),
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
		if errors.Is(err, pc.ErrTypeNotFound) {
			log.Warn("pc type not found", sl.Err(err))
			response.NotFound(w, "pc type not found")
			return
		}
		if errors.Is(err, pc.ErrPcNotFound) {
			log.Warn("pc not found")
			response.NotFound(w, "pc not found")
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

func (a *API) PcTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc.PcTypes"

		log := a.Log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req := &PcTypesRequest{
			limit:  r.URL.Query().Get("limit"),
			offset: r.URL.Query().Get("offset"),
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

		if req.offset == "" {
			req.offset = "0"
		}
		if req.limit == "" {
			req.limit = "0"
		}

		limit, err := strconv.ParseInt(req.limit, 10, 64)
		if err != nil {
			log.Error("failed to parse limit", sl.Err(err))
			response.Internal(w)
			return
		}

		offset, err := strconv.ParseInt(req.offset, 10, 64)
		if err != nil {
			log.Error("failed to parse offset", sl.Err(err))
			response.Internal(w)
			return
		}

		pcTypes, err := a.PcService.PcTypes(r.Context(), limit, offset)
		if err != nil {
			log.Error("failed to get pc types", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, pcTypes)
	}
}

func (a *API) PcType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pc.PcType"

		log := a.Log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req := &PcTypeRequest{
			typeId: chi.URLParam(r, "type_id"),
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

		typeId, err := strconv.ParseInt(req.typeId, 10, 64)
		if err != nil {
			log.Error("failed to parse type id", sl.Err(err))
			response.Internal(w)
			return
		}

		pcType, err := a.PcService.PcType(r.Context(), typeId)
		if errors.Is(err, pc.ErrTypeNotFound) {
			log.Warn("pc type not found", sl.Err(err))
			response.NotFound(w, "pc type not found")
			return
		}
		if err != nil {
			log.Error("failed to get pc type", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, pcType)
	}
}
