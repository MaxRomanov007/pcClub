package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrTokenMalformed = &Error{
		Code:    "TokenMalformed",
		Message: "token is malformed",
	}
	ErrTokenSignatureInvalid = &Error{
		Code:    "TokenSignatureInvalid",
		Message: "token signature is invalid",
	}
	ErrTokenExpired = &Error{
		Code:    "TokenExpired",
		Message: "token is expired",
	}
	ErrTokenInBlackList = &Error{
		Code:    "TokenInBlackList",
		Message: "token is in blacklist",
	}
	ErrUserNotFound = &Error{
		Code:    "UserNotFound",
		Message: "user not found",
	}
	ErrInvalidRefreshVersion = &Error{
		Code:    "InvalidRefreshVersion",
		Message: "invalid refresh version",
	}
)

func TokenError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		err = ErrTokenMalformed
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		err = ErrTokenSignatureInvalid
	case errors.Is(err, jwt.ErrTokenExpired), errors.Is(err, jwt.ErrTokenNotValidYet):
		err = ErrTokenExpired
	}

	return err
}
