package pcCLub

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/logger/sl"
	"server/internal/lib/response"
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

		if !a.validateRequest(w, *req, log) {
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

		req := &PcTypeRequest{
			typeId: chi.URLParam(r, "type_id"),
		}

		if !a.validateRequest(w, *req, log) {
			return
		}

		typeId, err := strconv.ParseInt(req.typeId, 10, 64)
		if err != nil {
			log.Error("failed to parse type id", sl.Err(err))
			response.Internal(w)
			return
		}

		pcType, err := a.PcService.PcType(r.Context(), typeId)
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
