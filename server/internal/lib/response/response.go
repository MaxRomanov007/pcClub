package response

import (
	"fmt"
	"github.com/go-playground/validator"
	"net/http"
	"server/internal/services/pcClub/auth"
	"server/internal/services/pcClub/components"
	"server/internal/services/pcClub/dish"
	"server/internal/services/pcClub/pc"
	"server/internal/services/pcClub/pcRoom"
	"server/internal/services/pcClub/user"
	"strings"
)

func Internal(w http.ResponseWriter) {
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func Unauthorized(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusUnauthorized)
}

func ValidationFailed(w http.ResponseWriter, errs validator.ValidationErrors) {
	http.Error(w, ValidationError(errs), http.StatusBadRequest)
}

func ValidationError(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		field := fieldName(err)

		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is required", field))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is over maximum", field))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is lover minimum", field))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is not valid", field))
		}
	}
	return strings.Join(errMsgs, "; ")
}

func fieldName(err validator.FieldError) string {
	return err.StructNamespace()[strings.Index(err.StructNamespace(), ".")+1:]
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
