package authAdmin

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"server/internal/lib/api/logger/sl"
	"server/internal/lib/api/request"
	"server/internal/lib/api/response"
	"server/internal/services/pcClub/auth"
	"server/internal/services/pcClub/user"
)

type AuthService interface {
	Access(
		ctx context.Context,
		accessToken string,
	) (uid int64, err error)
}

type UserService interface {
	IsAdmin(
		ctx context.Context,
		uid int64,
	) (err error)
}

func AuthAdmin(log *slog.Logger, s AuthService, u UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		const op = "middleware.auth.auth.authAdmin.AuthAdmin"
		fn := func(w http.ResponseWriter, r *http.Request) {
			log = log.With(
				slog.String("operation", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			access := request.AccessToken(w, r, log)
			if access == "" {
				log.Warn("no access token")
				response.Unauthorized(w, "no access token in header")
				return
			}

			uid, err := s.Access(r.Context(), access)
			if err != nil {
				var authError *auth.Error
				if ok := errors.As(err, &authError); ok {
					log.Warn("auth failed", sl.Err(err))
					response.AuthorizationFailed(w, authError)
					return
				}

				log.Error("failed to access token", sl.Err(err))
				response.Internal(w)
				return
			}

			if err := u.IsAdmin(r.Context(), uid); err != nil {
				if errors.Is(err, user.ErrAccessDenied) {
					log.Warn("access denied", sl.Err(err))
					response.Unauthorized(w, "access denied")
					return
				}

				log.Error("failed to check if user is admin", sl.Err(err))
				response.Internal(w)
				return
			}

			ctx := context.WithValue(r.Context(), "uid", uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
