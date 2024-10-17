package pcCLub

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	"server/internal/config"
	"server/internal/lib/api/cookie"
	"server/internal/lib/api/response"
	"server/internal/lib/logger/sl"
	"server/internal/models"
	"server/internal/services/pcClub/auth"
	"server/internal/services/pcClub/user"
	"strings"
)

type AuthService interface {
	Access(
		ctx context.Context,
		AccessToken string,
	) (uid int64, err error)

	Refresh(
		ctx context.Context,
		RefreshToken string,
	) (
		accessToken string,
		refreshToken string,
		err error,
	)

	Tokens(
		ctx context.Context,
		uid int64,
	) (
		accessToken string,
		refreshToken string,
		err error,
	)

	BanTokens(
		ctx context.Context,
		accessToken string,
		refreshToken string,
	) (uid int64, err error)
}

type UserService interface {
	SaveUser(
		ctx context.Context,
		email string,
		password string,
	) (uid int64, err error)

	Login(
		ctx context.Context,
		email string,
		password string,
	) (uid int64, err error)

	User(
		ctx context.Context,
		uid int64,
	) (user models.UserData, err error)

	UserByEmail(
		ctx context.Context,
		email string,
	) (user models.User, err error)

	DeleteUser(
		ctx context.Context,
		uid int64,
	) (err error)

	IsAdmin(
		ctx context.Context,
		uid int64,
	) (err error)
}

type PcService interface {
	Pcs(
		ctx context.Context,
		typeId int64,
		isAvailable bool,
	) (pcs []models.PcData, err error)

	PcTypes(
		ctx context.Context,
		limit int64,
		offset int64,
	) (pcs []models.PcTypeData, err error)

	PcType(
		ctx context.Context,
		typeId int64,
	) (pcType models.PcTypeData, err error)

	SavePcType(
		ctx context.Context,
		name string,
		description string,
		processor *models.ProcessorData,
		videoCard *models.VideoCardData,
		monitor *models.MonitorData,
		ram *models.RamData,
	) (err error)

	SavePc(
		ctx context.Context,
		typeId int64,
		roomId int64,
		row int,
		place int,
	) (err error)

	UpdatePcType(
		ctx context.Context,
		typeId int64,
		name string,
		description string,
		processor *models.ProcessorData,
		videoCard *models.VideoCardData,
		monitor *models.MonitorData,
		ram *models.RamData,
	) (err error)

	DeletePcType(
		ctx context.Context,
		typeId int64,
	) (err error)
}

type API struct {
	Log         *slog.Logger
	Cfg         *config.Config
	UserService UserService
	AuthService AuthService
	PcService   PcService
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	userService UserService,
	authService AuthService,
	pcService PcService,
) *API {
	return &API{
		Log:         log,
		Cfg:         cfg,
		UserService: userService,
		AuthService: authService,
		PcService:   pcService,
	}
}

func (a *API) decodeAndValidateRequest(w http.ResponseWriter, r *http.Request, log *slog.Logger, req interface{}) bool {
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		response.Internal(w)
		return false
	}

	if err := validator.New().Struct(req); err != nil {
		var validError validator.ValidationErrors
		if ok := errors.As(err, &validError); ok {
			log.Warn("invalid request", sl.Err(err))
			response.ValidationFailed(w, validError)
			return false
		}

		log.Error("failed to validate request", sl.Err(err))
		response.Internal(w)
		return false
	}

	return true
}

func (a *API) setRefreshCookie(w http.ResponseWriter, refreshToken string) {
	cookie.Set(
		w,
		a.Cfg.Auth.Refresh.CookieName,
		refreshToken,
		a.Cfg.Auth.Path,
		a.Cfg.Auth.Refresh.TTL,
	)
}

func (a *API) getRefreshToken(w http.ResponseWriter, r *http.Request, log *slog.Logger) string {
	refresh, err := r.Cookie(a.Cfg.Auth.Refresh.CookieName)
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

func (a *API) getAccessToken(w http.ResponseWriter, r *http.Request, log *slog.Logger) string {
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

func (a *API) authorizeRequest(w http.ResponseWriter, r *http.Request, log *slog.Logger) int64 {
	access := a.getAccessToken(w, r, log)
	if access == "" {
		return 0
	}

	uid, err := a.AuthService.Access(r.Context(), access)
	if err != nil {
		var authError *auth.Error
		if ok := errors.As(err, &authError); ok {
			log.Warn("auth failed", sl.Err(err))
			response.AuthorizationFailed(w, authError)
			return 0
		}

		log.Error("failed to access token", sl.Err(err))
		response.Internal(w)
		return 0
	}

	return uid
}

func (a *API) authorizeAdmin(w http.ResponseWriter, r *http.Request, log *slog.Logger) bool {
	uid := a.authorizeRequest(w, r, log)
	if uid == 0 {
		return false
	}

	if err := a.UserService.IsAdmin(r.Context(), uid); err != nil {
		if errors.Is(err, user.ErrAccessDenied) {
			log.Warn("access denied", sl.Err(err))
			response.Unauthorized(w, "access denied")
			return false
		}

		log.Error("failed to check if user is admin", sl.Err(err))
		response.Internal(w)
		return false
	}

	return true
}

func (a *API) log(op string, r *http.Request) *slog.Logger {
	return a.Log.With(
		slog.String("operation", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
}
