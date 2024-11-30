package request

import (
	"errors"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"server/internal/lib/api/logger/sl"
	"server/internal/lib/api/response"
)

func ValidateRequest[T any](w http.ResponseWriter, req any, log *slog.Logger) bool {
	if err := validator.New().Struct(req); err != nil {
		var validError validator.ValidationErrors
		if ok := errors.As(err, &validError); ok {
			log.Warn("invalid request", sl.Err(err))
			response.ValidationFailed[T](w, validError)
			return false
		}

		log.Error("failed to validate request", sl.Err(err))
		response.Internal(w)
		return false
	}
	return true
}
