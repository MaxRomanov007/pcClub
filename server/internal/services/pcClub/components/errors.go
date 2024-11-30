package components

import (
	"errors"
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
	ErrNotFoundCode           = "NotFound"
	ErrAlreadyExistsCode      = "AlreadyExists"
	ErrReferenceNotExistsCode = "ReferenceNotExists"
)

var (
	ErrNotFound = &Error{
		Code:    ErrNotFoundCode,
		Message: "not found",
	}
	ErrAlreadyExists = &Error{
		Code:    ErrAlreadyExistsCode,
		Message: "already exists",
	}
	ErrReferenceNotExists = &Error{
		Code:    ErrReferenceNotExistsCode,
		Message: "reference not exists",
	}
)

func HandleStorageError(err error) error {
	var ssmsErr *gorm.Error
	if !errors.As(err, &ssmsErr) {
		return errors2.WithMessage(err, "unknown error")
	}
	switch ssmsErr.Code {
	case gorm.ErrNotFoundCode:
		err = ErrNotFound
	case gorm.ErrAlreadyExistsCode:
		err = ErrAlreadyExists
	case gorm.ErrReferenceNotExistsCode:
		err = ErrReferenceNotExists
	default:
		err = errors2.WithMessage(ssmsErr, "unknown mssql error")
	}

	return err
}
