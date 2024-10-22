package ssms

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
