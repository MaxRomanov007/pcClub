package auth

import (
	"context"
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
