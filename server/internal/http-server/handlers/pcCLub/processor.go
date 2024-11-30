package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/api/logger/sl"
	"server/internal/lib/api/request"
	"server/internal/lib/api/response"
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

		req, ok := request.DecodeAndValidateGETRequest[ProcessorsRequest](w, r, log)
		if !ok {
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

		req, ok := request.DecodeAndValidateJSONRequest[SaveProcessorProducerRequest](w, r, log)
		if !ok {
			return
		}

		producer := models.ProcessorProducer{
			Name: req.Name,
		}
		if _, err := a.ComponentsService.Processor.SaveProcessorProducer(r.Context(), &producer); err != nil {
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

		render.JSON(w, r, producer)
	}
}

func (a *API) SaveProcessor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.processor.SaveProcessor"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[SaveProcessorRequest](w, r, log)
		if !ok {
			return
		}

		processor := models.Processor{
			ProcessorProducerID: req.ProducerId,
			Model:               req.Model,
		}
		if _, err := a.ComponentsService.Processor.SaveProcessor(r.Context(), &processor); err != nil {
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

		render.JSON(w, r, processor)
	}
}

func (a *API) DeleteProcessorProducer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.processor.DeleteProcessorProducer"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[DeleteProcessorProducerRequest](w, r, log)
		if !ok {
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

		req, ok := request.DecodeAndValidateJSONRequest[DeleteProcessorRequest](w, r, log)
		if !ok {
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
