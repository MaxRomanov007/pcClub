package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"server/internal/config"
	"time"
)

type redisProvider interface {
	StringValue(
		ctx context.Context,
		key string,
	) (value string, err error)
}

type redisOwner interface {
	SetStringWithCustomTTL(
		ctx context.Context,
		key string,
		value string,
		ttl time.Duration,
	) (err error)
}

type versionProvider interface {
	RefreshVersion(
		ctx context.Context,
		uid int64,
	) (version int64, err error)
}

type versionOwner interface {
	UpdateRefreshVersion(
		ctx context.Context,
		uid int64,
		version int64,
	) (err error)
}

type Service struct {
	cfg             *config.AuthConfig
	redisOwner      redisOwner
	redisProvider   redisProvider
	versionProvider versionProvider
	versionOwner    versionOwner
}

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

func New(
	cfg *config.AuthConfig,
	redisOwner redisOwner,
	redisProvider redisProvider,
	VersionOwner versionOwner,
	VersionProvider versionProvider,
) *Service {
	return &Service{
		cfg:             cfg,
		redisOwner:      redisOwner,
		redisProvider:   redisProvider,
		versionOwner:    VersionOwner,
		versionProvider: VersionProvider,
	}
}

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
