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

		req, ok := request.DecodeAndValidateGETRequest[VideoCardsRequest](w, r, log)
		if !ok {
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

		req, ok := request.DecodeAndValidateJSONRequest[SaveVideoCardProducerRequest](w, r, log)
		if !ok {
			return
		}

		producer := models.VideoCardProducer{
			Name: req.Name,
		}
		if _, err := a.ComponentsService.VideoCard.SaveVideoCardProducer(r.Context(), &producer); err != nil {
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

		render.JSON(w, r, producer)
	}
}

func (a *API) SaveVideoCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.videoCard.SaveVideoCard"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[SaveVideoCardRequest](w, r, log)
		if !ok {
			return
		}

		card := models.VideoCard{
			VideoCardProducerID: req.ProducerId,
			Model:               req.Model,
		}
		if _, err := a.ComponentsService.VideoCard.SaveVideoCard(r.Context(), &card); err != nil {
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

		render.JSON(w, r, card)
	}
}

func (a *API) DeleteVideoCardProducer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pcClub.videoCard.DeleteVideoCardProducer"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[DeleteVideoCardProducerRequest](w, r, log)
		if !ok {
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

		req, ok := request.DecodeAndValidateJSONRequest[DeleteVideoCardRequest](w, r, log)
		if !ok {
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
