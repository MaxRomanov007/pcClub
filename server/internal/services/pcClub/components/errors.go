package components

import "server/internal/storage/ssms"

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

func HandleStorageError(ssmsErr *ssms.Error) *Error {
	switch ssmsErr.Code {
	case ssms.ErrNotFoundCode:
		return ErrNotFound
	case ssms.ErrAlreadyExistsCode:
		return ErrAlreadyExists
	case ssms.ErrReferenceNotExistsCode:
		return ErrReferenceNotExists
	default:
		return &Error{
			Message: ssmsErr.Message,
		}
	}
}
