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

type ProcessorsRequest struct {
	ProducerId int64 `get:"producer-id" validate:"required,min=1"`
}

type SaveProcessorProducerRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type SaveProcessorRequest struct {
	ProducerId int64  `json:"producer_id" validate:"required,min=1"`
	Model      string `json:"model" validate:"required,max=255"`
}

type DeleteProcessorProducerRequest struct {
	ProducerId int64 `json:"producer_id" validate:"required,min=1"`
}

type DeleteProcessorRequest struct {
	ProcessorId int64 `json:"processor_id" validate:"required,min=1"`
}

func (a *API) ProcessorProducers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.processor.ProcessorProducers"

		log := a.log(op, r)

		producers, err := a.ComponentsService.Processor.ProcessorProducers(r.Context())
		if err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to get processor producers", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, producers)
	}
}

func (a *API) Processors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.processor.Processors"

		log := a.log(op, r)

		var req ProcessorsRequest
		if !a.decodeAndValidateGETRequest(w, r, log, &req) {
			return
		}

		processors, err := a.ComponentsService.Processor.Processors(r.Context(), req.ProducerId)
		if err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to get processors", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, processors)
	}
}

func (a *API) SaveProcessorProducer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.processor.SaveProcessorProducer"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SaveProcessorProducerRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.Processor.SaveProcessorProducer(r.Context(), req.Name); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to save processor producer", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) SaveProcessor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.processor.SaveProcessor"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SaveProcessorRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.Processor.SaveProcessor(r.Context(), models.Processor{
			ProcessorProducerID: req.ProducerId,
			Model:               req.Model,
		}); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to save processor", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) DeleteProcessorProducer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.processor.DeleteProcessorProducer"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req DeleteProcessorProducerRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.Processor.DeleteProcessorProducer(r.Context(), req.ProducerId); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to delete processor producer", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) DeleteProcessor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.processor.DeleteProcessor"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req DeleteProcessorRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.Processor.DeleteProcessor(r.Context(), req.ProcessorId); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to delete processor", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}
