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

type VideoCardsRequest struct {
	ProducerId int64 `get:"producer-id" validate:"required,min=1"`
}

type SaveVideoCardProducerRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type SaveVideoCardRequest struct {
	ProducerId int64  `json:"producer_id" validate:"required,min=1"`
	Model      string `json:"model" validate:"required,max=255"`
}

type DeleteVideoCardProducerRequest struct {
	ProducerId int64 `json:"producer_id" validate:"required,min=1"`
}

type DeleteVideoCardRequest struct {
	VideoCardId int64 `json:"videoCard_id" validate:"required,min=1"`
}

func (a *API) VideoCardProducers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.videoCard.VideoCardProducers"

		log := a.log(op, r)

		producers, err := a.ComponentsService.VideoCard.VideoCardProducers(r.Context())
		if err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to get videoCard producers", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, producers)
	}
}

func (a *API) VideoCards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.videoCard.VideoCards"

		log := a.log(op, r)

		var req VideoCardsRequest
		if !a.decodeAndValidateGETRequest(w, r, log, &req) {
			return
		}

		videoCards, err := a.ComponentsService.VideoCard.VideoCards(r.Context(), req.ProducerId)
		if err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to get videoCards", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, videoCards)
	}
}

func (a *API) SaveVideoCardProducer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.videoCard.SaveVideoCardProducer"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SaveVideoCardProducerRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.VideoCard.SaveVideoCardProducer(r.Context(), req.Name); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to save videoCard producer", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) SaveVideoCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.videoCard.SaveVideoCard"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req SaveVideoCardRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.VideoCard.SaveVideoCard(r.Context(), models.VideoCard{
			VideoCardProducerID: req.ProducerId,
			Model:               req.Model,
		}); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to save videoCard", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) DeleteVideoCardProducer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.videoCard.DeleteVideoCardProducer"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req DeleteVideoCardProducerRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.VideoCard.DeleteVideoCardProducer(r.Context(), req.ProducerId); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to delete videoCard producer", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}

func (a *API) DeleteVideoCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.videoCard.DeleteVideoCard"

		log := a.log(op, r)

		if !a.authorizeAdmin(w, r, log) {
			return
		}

		var req DeleteVideoCardRequest
		if !a.decodeAndValidateJSONRequest(w, r, log, &req) {
			return
		}

		if err := a.ComponentsService.VideoCard.DeleteVideoCard(r.Context(), req.VideoCardId); err != nil {
			var serviceErr *components.Error
			if errors.As(err, &serviceErr) {
				log.Warn("component error", sl.Err(err))
				response.ComponentsError(w, serviceErr)
				return
			}
			log.Error("failed to delete videoCard", sl.Err(err))
			response.Internal(w)
			return
		}
	}
}
