package user

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
	ErrUserAlreadyExists = &Error{
		Code:    ErrUserAlreadyExistsCode,
		Message: "user already exists",
	}
	ErrAccessDenied = &Error{
		Code:    ErrAccessDeniedCode,
		Message: "access denied",
	}
	ErrUserNotFound = &Error{
		Code:    ErrUserNotFoundCode,
		Message: "user not found",
	}
)
