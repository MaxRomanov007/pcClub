package pcCLub

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"net/http"
	"server/internal/lib/api/response"
	"server/internal/lib/logger/sl"
	"server/internal/services/pcClub/pc"
	"strconv"
)

type PcTypesRequest struct {
	limit  string `validate:"omitempty,number"`
	offset string `validate:"omitempty,number"`
}

type PcTypeRequest struct {
	typeId string `validate:"require,number"`
}

func (a *API) PcTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.pcType.PcTypes"

		log := a.log(op, r)

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
		const op = "handlers.pcClub.pcType.PcType"

		log := a.log(op, r)

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
		if errors.Is(err, pc.ErrNotFound) {
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
