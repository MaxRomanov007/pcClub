package response

import (
	"github.com/go-playground/validator"
	"net/http"
	"server/internal/config"
	validator2 "server/internal/lib/api/validator"
	"server/internal/lib/cookie"
	"server/internal/services/pcClub/auth"
	"server/internal/services/pcClub/components"
	"server/internal/services/pcClub/dish"
	"server/internal/services/pcClub/orderPc"
	"server/internal/services/pcClub/pc"
	"server/internal/services/pcClub/pcRoom"
	"server/internal/services/pcClub/user"
)

func Internal(w http.ResponseWriter) {
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func Unauthorized(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusUnauthorized)
}

func ValidationFailed[T any](w http.ResponseWriter, errs validator.ValidationErrors) {
	http.Error(w, validator2.ValidationError[T](errs), http.StatusBadRequest)
}

func AuthorizationFailed(w http.ResponseWriter, err *auth.Error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}

func PcError(w http.ResponseWriter, err *pc.Error) {
	switch err.Code {
	case pc.ErrNotFoundCode:
		http.Error(w, err.Error(), http.StatusNotFound)
	case pc.ErrAlreadyExistsCode, pc.ErrConstraintCode, pc.ErrReferenceNotExistsCode:
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		Internal(w)
	}
}

func OrderError(w http.ResponseWriter, err *orderPc.Error) {
	switch err.Code {
	case orderPc.ErrNotFoundCode:
		http.Error(w, err.Error(), http.StatusNotFound)
	case orderPc.ErrAlreadyExistsCode, orderPc.ErrConstraintCode, orderPc.ErrReferenceNotExistsCode:
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		Internal(w)
	}
}

func UserError(w http.ResponseWriter, err *user.Error) {
	switch err.Code {
	case user.ErrUserNotFoundCode:
		http.Error(w, err.Error(), http.StatusNotFound)
	case user.ErrUserAlreadyExistsCode:
		http.Error(w, err.Error(), http.StatusConflict)
	case user.ErrAccessDeniedCode, user.ErrInvalidCredentialsCode:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	default:
		Internal(w)
	}
}

func PcRoomError(w http.ResponseWriter, err *pcRoom.Error) {
	switch err.Code {
	case pcRoom.ErrNotFoundCode:
		http.Error(w, err.Error(), http.StatusNotFound)
	case pcRoom.ErrAlreadyExistsCode, pcRoom.ErrReferenceNotExistsCode:
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		Internal(w)
	}
}

func ComponentsError(w http.ResponseWriter, err *components.Error) {
	switch err.Code {
	case components.ErrNotFoundCode:
		http.Error(w, "not found", http.StatusNotFound)
	case components.ErrAlreadyExistsCode:
		http.Error(w, "already exists", http.StatusConflict)
	case components.ErrReferenceNotExistsCode:
		http.Error(w, "reference doesnt exists", http.StatusConflict)
	default:
		Internal(w)
	}
}

func DishError(w http.ResponseWriter, err *dish.Error) {
	switch err.Code {
	case components.ErrNotFoundCode:
		http.Error(w, "not found", http.StatusNotFound)
	case components.ErrAlreadyExistsCode:
		http.Error(w, "already exists", http.StatusConflict)
	case components.ErrReferenceNotExistsCode:
		http.Error(w, "reference doesnt exists", http.StatusConflict)
	default:
		Internal(w)
	}
}

func SetRefreshCookie(w http.ResponseWriter, cfg *config.AuthConfig, refreshToken string) {
	cookie.Set(
		w,
		cfg.Refresh.CookieName,
		refreshToken,
		cfg.UrlPath,
		cfg.Refresh.TTL,
	)
}
