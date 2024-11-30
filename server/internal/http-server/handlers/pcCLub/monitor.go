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

type MonitorsRequest struct {
	ProducerId int64 `get:"producer-id" validate:"required,min=1"`
}

type SaveMonitorProducerRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type SaveMonitorRequest struct {
	ProducerId int64  `json:"producer_id" validate:"required,min=1"`
	Model      string `json:"model" validate:"required,max=255"`
}

type DeleteMonitorProducerRequest struct {
	ProducerId int64 `json:"producer_id" validate:"required,min=1"`
}

type DeleteMonitorRequest struct {
	MonitorId int64 `json:"monitor_id" validate:"required,min=1"`
}

func (a *API) MonitorProducers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.monitor.MonitorProducers"

		log := a.log(op, r)

		producers, err := a.ComponentsService.Monitor.MonitorProducers(r.Context())
		if err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to get monitor producers", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, producers)
	}
}

func (a *API) Monitors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.monitor.Monitors"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateGETRequest[MonitorsRequest](w, r, log)
		if !ok {
			return
		}

		monitors, err := a.ComponentsService.Monitor.Monitors(r.Context(), req.ProducerId)
		if err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to get monitors", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, monitors)
	}
}

func (a *API) SaveMonitorProducer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.monitor.SaveMonitorProducer"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[SaveMonitorProducerRequest](w, r, log)
		if !ok {
			return
		}

		producer := models.MonitorProducer{
			Name: req.Name,
		}
		if _, err := a.ComponentsService.Monitor.SaveMonitorProducer(r.Context(), &producer); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to save monitor producer", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, producer)
	}
}

func (a *API) SaveMonitor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.monitor.SaveMonitor"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[SaveMonitorRequest](w, r, log)
		if !ok {
			return
		}

		monitor := models.Monitor{
			MonitorProducerID: req.ProducerId,
			Model:             req.Model,
		}
		if _, err := a.ComponentsService.Monitor.SaveMonitor(r.Context(), &monitor); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to save monitor", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, monitor)
	}
}

func (a *API) DeleteMonitorProducer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.monitor.DeleteMonitorProducer"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[DeleteMonitorProducerRequest](w, r, log)
		if !ok {
			return
		}

		if err := a.ComponentsService.Monitor.DeleteMonitorProducer(r.Context(), req.ProducerId); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to delete monitor producer", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) DeleteMonitor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.monitor.DeleteMonitor"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[DeleteMonitorRequest](w, r, log)
		if !ok {
			return
		}

		if err := a.ComponentsService.Monitor.DeleteMonitor(r.Context(), req.MonitorId); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to delete monitor", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}
