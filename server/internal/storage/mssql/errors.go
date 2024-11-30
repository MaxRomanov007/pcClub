package mssql

import (
	"errors"
	"gorm.io/gorm"
	errors2 "server/internal/lib/errors"
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
	ErrCheckFailedCode        = "CheckFailed"
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

func errorByResult(res *gorm.DB) error {
	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
		return ErrAlreadyExists
	}
	if errors.Is(res.Error, gorm.ErrForeignKeyViolated) {
		return ErrReferenceNotExists
	}
	if errors.Is(res.Error, gorm.ErrRecordNotFound) ||
		res.Error == nil && res.RowsAffected == 0 {
		return ErrNotFound
	}
	return errors2.WithMessage(res.Error, "unknown error")
}
