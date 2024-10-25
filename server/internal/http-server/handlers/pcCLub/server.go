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
	"server/internal/lib/cookie"
	"server/internal/lib/logger/sl"
	"server/internal/lib/request/urlGet"
	"server/internal/lib/response"
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

type PcTypeService interface {
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

type PcService interface {
	Pcs(
		ctx context.Context,
		typeId int64,
		isAvailable bool,
	) (pcs []models.PcData, err error)

	SavePc(
		ctx context.Context,
		typeId int64,
		roomId int64,
		row int,
		place int,
		description string,
	) (err error)

	UpdatePc(
		ctx context.Context,
		pcId int64,
		typeId int64,
		roomId int64,
		statusId int64,
		row int,
		place int,
		description string,
	) (err error)

	DeletePc(
		ctx context.Context,
		pcId int64,
	) (err error)
}

type PcRoomService interface {
	PcRoom(
		ctx context.Context,
		pcId int64,
	) (pcRoom models.PcRoom, err error)

	SavePcRoom(
		ctx context.Context,
		pcRoom models.PcRoom,
	) (err error)

	UpdatePcRoom(
		ctx context.Context,
		pcRoom models.PcRoom,
	) (err error)

	DeletePcRoom(
		ctx context.Context,
		pcId int64,
	) (err error)
}

type ProcessorService interface {
	ProcessorProducers(
		ctx context.Context,
	) (producers []models.ProcessorProducer, err error)

	Processors(
		ctx context.Context,
		producerId int64,
	) (processors []models.Processor, err error)

	SaveProcessorProducer(
		ctx context.Context,
		name string,
	) (err error)

	SaveProcessor(
		ctx context.Context,
		processor models.Processor,
	) (err error)

	DeleteProcessorProducer(
		ctx context.Context,
		producerId int64,
	) (err error)

	DeleteProcessor(
		ctx context.Context,
		processorId int64,
	) (err error)
}

type MonitorService interface {
	MonitorProducers(
		ctx context.Context,
	) (producers []models.MonitorProducer, err error)

	Monitors(
		ctx context.Context,
		producerId int64,
	) (monitors []models.Monitor, err error)

	SaveMonitorProducer(
		ctx context.Context,
		name string,
	) (err error)

	SaveMonitor(
		ctx context.Context,
		monitor models.Monitor,
	) (err error)

	DeleteMonitorProducer(
		ctx context.Context,
		producerId int64,
	) (err error)

	DeleteMonitor(
		ctx context.Context,
		monitorId int64,
	) (err error)
}

type VideoCardService interface {
	VideoCardProducers(
		ctx context.Context,
	) (producers []models.VideoCardProducer, err error)

	VideoCards(
		ctx context.Context,
		producerId int64,
	) (videoCards []models.VideoCard, err error)

	SaveVideoCardProducer(
		ctx context.Context,
		name string,
	) (err error)

	SaveVideoCard(
		ctx context.Context,
		videoCard models.VideoCard,
	) (err error)

	DeleteVideoCardProducer(
		ctx context.Context,
		producerId int64,
	) (err error)

	DeleteVideoCard(
		ctx context.Context,
		videoCardId int64,
	) (err error)
}

type RamService interface {
	RamTypes(
		ctx context.Context,
	) (producers []models.RamType, err error)

	Rams(
		ctx context.Context,
		typeId int64,
	) (rams []models.Ram, err error)

	SaveRamType(
		ctx context.Context,
		name string,
	) (err error)

	SaveRam(
		ctx context.Context,
		ram models.Ram,
	) (err error)

	DeleteRamType(
		ctx context.Context,
		typeId int64,
	) (err error)

	DeleteRam(
		ctx context.Context,
		ramId int64,
	) (err error)
}

type ComponentsService struct {
	Processor ProcessorService
	Monitor   MonitorService
	VideoCard VideoCardService
	Ram       RamService
}

type API struct {
	Log               *slog.Logger
	Cfg               *config.Config
	UserService       UserService
	AuthService       AuthService
	PcTypeService     PcTypeService
	PcService         PcService
	PcRoomService     PcRoomService
	ComponentsService ComponentsService
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	userService UserService,
	authService AuthService,
	pcTypeService PcTypeService,
	pcService PcService,
	pcRoomService PcRoomService,
	componentsService ComponentsService,
) *API {
	return &API{
		Log:               log,
		Cfg:               cfg,
		UserService:       userService,
		AuthService:       authService,
		PcTypeService:     pcTypeService,
		PcService:         pcService,
		PcRoomService:     pcRoomService,
		ComponentsService: componentsService,
	}
}

func (a *API) decodeAndValidateJSONRequest(
	w http.ResponseWriter,
	r *http.Request,
	log *slog.Logger,
	req any,
) bool {
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		response.Internal(w)
		return false
	}

	if !a.validateRequest(w, req, log) {
		return false
	}

	return true
}

func (a *API) decodeAndValidateGETRequest(
	w http.ResponseWriter,
	r *http.Request,
	log *slog.Logger,
	req any,
) bool {
	if err := urlGet.Decode(r, req); err != nil {
		log.Error("failed to decode get request", sl.Err(err))
		response.Internal(w)
		return false
	}
	if !a.validateRequest(w, req, log) {
		return false
	}
	return true
}

func (a *API) validateRequest(w http.ResponseWriter, req any, log *slog.Logger) bool {
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
