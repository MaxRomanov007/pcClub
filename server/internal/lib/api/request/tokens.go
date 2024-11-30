package request

import (
	"errors"
	"log/slog"
	"net/http"
	"server/internal/lib/api/logger/sl"
	"server/internal/lib/api/response"
	"strings"
)

func RefreshToken(w http.ResponseWriter, r *http.Request, log *slog.Logger, cookieName string) string {
	refresh, err := r.Cookie(cookieName)
	if errors.Is(err, http.ErrNoCookie) {
		log.Warn("no refresh cookie in request", sl.Err(err))
		response.Unauthorized(w, "no refresh cookie in request")
		return ""
	}
	if err != nil {
		log.Error("failed to get cookie", sl.Err(err))
		response.Internal(w)
		return ""
	}

	return refresh.Value
}

func AccessToken(w http.ResponseWriter, r *http.Request, log *slog.Logger) string {
	access := r.Header.Get("Authorization")
	if access == "" {
		log.Warn("no authorization header in request")
		response.Unauthorized(w, "no authorization header in request")
		return ""
	}

	access, isBearer := strings.CutPrefix(access, "Bearer ")
	if !isBearer {
		log.Warn("auth header is not bearer token")
		response.Unauthorized(w, "auth header is not bearer token")
		return ""
	}

	return access
}
