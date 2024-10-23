package pcType

import "server/internal/storage/ssms"

type Error struct {
	Code    string
	Message string
	Desc    string
}

func (e *Error) Error() string {
	message := e.Message
	if e.Desc != "" {
		message += ": " + e.Desc
	}
	return message
}

const (
	ErrNotFoundCode           = "NotFound"
	ErrAlreadyExistsCode      = "AlreadyExists"
	ErrConstraintCode         = "Constraint"
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
	ErrConstraint = &Error{
		Code:    ErrConstraintCode,
		Message: "constraint failure",
	}
	ErrReferenceNotExists = &Error{
		Code:    ErrReferenceNotExistsCode,
		Message: "reference not exists",
	}
)

func (e *Error) WithDesc(desc string) *Error {
	e.Desc = desc
	return e
}

func handleStorageError(ssmsErr *ssms.Error) *Error {
	switch ssmsErr.Code {
	case ssms.ErrNotFoundCode:
		return ErrNotFound
	case ssms.ErrAlreadyExistsCode:
		return ErrAlreadyExists
	case ssms.ErrReferenceNotExistsCode:
		return ErrReferenceNotExists
	case ssms.ErrCheckFailedCode:
		return ErrConstraint
	default:
		return &Error{
			Message: ssmsErr.Message,
		}
	}
}
