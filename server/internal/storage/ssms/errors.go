package ssms

import (
	"database/sql"
	"errors"
	sqld "github.com/denisenkom/go-mssqldb"
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
	ErrTooLongCode            = "TooLong"
	ErrNullPointerCode        = "NullPointer"
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
	ErrTooLong = &Error{
		Code:    ErrTooLongCode,
		Message: "too long",
	}
	ErrNullPointer = &Error{
		Code:    ErrNullPointerCode,
		Message: "null pointer",
	}
	ErrCheckFailed = &Error{
		Code:    ErrCheckFailedCode,
		Message: "check failed",
	}
)

func handleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	var driverError sqld.Error
	if ok := errors.As(err, &driverError); ok {
		switch driverError.Number {
		case 2627: // Нарушение уникального индекса
			return ErrAlreadyExists
		case 547: // Нарушение внешнего ключа
			return ErrReferenceNotExists
		case 8152: // Строка слишком длинная для столбца
			return ErrTooLong
		case 515:
			return ErrNullPointer
		case 54332:
			return ErrCheckFailed
		default:
			return err
		}
	}
	return err
}
