package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	errors2 "server/internal/lib/errors"
	gorm "server/internal/storage/mssql"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

const (
	ErrInvalidCredentialsCode = "InvalidCredentials"
	ErrUserAlreadyExistsCode  = "UserAlreadyExists"
	ErrAccessDeniedCode       = "AccessDenied"
	ErrUserNotFoundCode       = "UserNotFound"
)

var (
	ErrInvalidCredentials = &Error{
		Code:    ErrInvalidCredentialsCode,
		Message: "credentials are not valid",
	}
	ErrAlreadyExists = &Error{
		Code:    ErrUserAlreadyExistsCode,
		Message: "user already exists",
	}
	ErrAccessDenied = &Error{
		Code:    ErrAccessDeniedCode,
		Message: "access denied",
	}
	ErrNotFound = &Error{
		Code:    ErrUserNotFoundCode,
		Message: "user not found",
	}
)

func HandleStorageError(err error) error {
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidCredentials
	}
	var ssmsErr *gorm.Error
	if !errors.As(err, &ssmsErr) {
		return errors2.WithMessage(err, "unknown error")
	}
	switch ssmsErr.Code {
	case gorm.ErrNotFoundCode:
		err = ErrNotFound
	case gorm.ErrAlreadyExistsCode:
		err = ErrAlreadyExists
	default:
		err = errors2.WithMessage(ssmsErr, "unknown mssql error")
	}

	return err
}
