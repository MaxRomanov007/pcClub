package response

import (
	"fmt"
	"github.com/go-playground/validator"
	"net/http"
	"server/internal/services/pcClub/auth"
	"strings"
)

func Internal(w http.ResponseWriter) {
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func NotFound(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func AlreadyExists(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusConflict)
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
