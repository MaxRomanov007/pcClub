package pcCLub

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"server/internal/lib/api/logger/sl"
	"server/internal/lib/api/request"
	"server/internal/lib/api/response"
	"server/internal/services/pcClub/auth"
	"server/internal/services/pcClub/user"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,min=3,max=32,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
type RegisterResponse struct {
	Access string `json:"access"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,min=3,max=32,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
type LoginResponse struct {
	Access string `json:"access_token"`
}

type RefreshResponse struct {
	Access string `json:"access_token"`
}

func (a *API) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.SaveUser"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[RegisterRequest](w, r, log)
		if !ok {
			return
		}

		id, err := a.UserService.SaveUser(r.Context(), req.Email, req.Password)
		if err != nil {
			var userErr *user.Error
			if ok := errors.As(err, &userErr); ok {
				log.Warn("user error", sl.Err(err))
				response.UserError(w, userErr)
				return
			}
			log.Error("failed to register user", sl.Err(err))
			response.Internal(w)
			return
		}

		access, refresh, err := a.AuthService.Tokens(r.Context(), id)
		if err != nil {
			var authError *auth.Error
			if ok := errors.As(err, &authError); ok {
				log.Warn("auth failed", sl.Err(err))
				response.AuthorizationFailed(w, authError)
				return
			}

			log.Error("failed to get tokens", sl.Err(err))
			response.Internal(w)
			return
		}

		response.SetRefreshCookie(w, a.Cfg.Auth, refresh)

		render.JSON(w, r, RegisterResponse{
			Access: access,
		})
	}
}

func (a *API) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.SaveUser"

		log := a.log(op, r)

		req, ok := request.DecodeAndValidateJSONRequest[LoginRequest](w, r, log)
		if !ok {
			return
		}

		id, err := a.UserService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			var userErr *user.Error
			if ok := errors.As(err, &userErr); ok {
				log.Warn("user error", sl.Err(err))
				response.UserError(w, userErr)
				return
			}
			log.Error("failed to login user", sl.Err(err))
			response.Internal(w)
			return
		}

		access, refresh, err := a.AuthService.Tokens(r.Context(), id)
		if err != nil {
			var authError *auth.Error
			if ok := errors.As(err, &authError); ok {
				log.Warn("auth failed", sl.Err(err))
				response.AuthorizationFailed(w, authError)
				return
			}

			log.Error("failed to get tokens", sl.Err(err))
			response.Internal(w)
			return
		}

		response.SetRefreshCookie(w, a.Cfg.Auth, refresh)

		render.JSON(w, r, LoginResponse{
			Access: access,
		})
	}
}

func (a *API) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Refresh"

		log := a.log(op, r)

		refresh := request.RefreshToken(w, r, log, a.Cfg.Auth.Refresh.CookieName)
		if refresh == "" {
			return
		}

		access, refresh, err := a.AuthService.Refresh(r.Context(), refresh)
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

		response.SetRefreshCookie(w, a.Cfg.Auth, refresh)

		render.JSON(w, r, LoginResponse{
			Access: access,
		})
	}
}

func (a *API) User() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Refresh"

		log := a.log(op, r)

		uid := request.MustUID(r)

		userData, err := a.UserService.User(r.Context(), uid)
		if err != nil {
			var userErr *user.Error
			if ok := errors.As(err, &userErr); ok {
				log.Warn("user error", sl.Err(err))
				response.UserError(w, userErr)
				return
			}
			log.Error("failed to get user", sl.Err(err))
			response.Internal(w)
			return
		}

		render.JSON(w, r, userData)
	}
}

func (a *API) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Logout"

		log := a.log(op, r)

		access := request.AccessToken(w, r, log)
		if access == "" {
			return
		}

		refreshToken := request.RefreshToken(w, r, log, a.Cfg.Auth.Refresh.CookieName)
		if refreshToken == "" {
			return
		}

		_, err := a.AuthService.BanTokens(r.Context(), access, refreshToken)
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
	}
}
