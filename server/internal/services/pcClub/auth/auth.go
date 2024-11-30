package auth

import (
	"context"
	"errors"
	"fmt"
	"server/internal/config"
	errors2 "server/internal/lib/errors"
	"server/internal/lib/jwt"
	"server/internal/storage/mssql"
	"server/internal/storage/redis"
	"strconv"
)

// generateTokens generates access and refresh tokens
// (1 output is access, 2 is refresh)
func generateTokens(
	uid int64,
	refreshVersion int64,
	cfg *config.AuthConfig,
) (access string, refresh string, err error) {
	access, err = jwt.NewAccessToken(
		uid,
		cfg.Access.Secret,
		cfg.Access.TTL,
	)
	if err != nil {
		return "", "", err
	}

	refresh, err = jwt.NewRefreshToken(
		uid,
		refreshVersion,
		cfg.Refresh.Secret,
		cfg.Refresh.TTL,
	)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) Access(
	ctx context.Context,
	accessToken string,
) (int64, error) {
	const op = "services.pcClub.auth.Access"

	claims, err := jwt.ParseToken(accessToken, s.cfg.Access.Secret)
	if err != nil {
		return 0, errors2.WithMessage(TokenError(err), op, "failed to parse token")
	}

	banedTokenExpirationTimeString, err := s.redisProvider.StringValue(
		ctx,
		fmt.Sprintf("%s:%d", AccessRedisBlackListName, claims.UID),
	)
	if errors.Is(err, redis.ErrNotFound) {
		return claims.UID, nil
	}
	if err != nil {
		return 0, errors2.WithMessage(err, op, "failed to find token in redis")
	}

	var banedTokenExpirationTime int64
	banedTokenExpirationTime, err = strconv.ParseInt(banedTokenExpirationTimeString, 10, 64)
	if err != nil {
		return 0, errors2.WithMessage(err, op, "failed to parse banned token expiration time")
	}

	if banedTokenExpirationTime >= claims.ExpiresAt.Unix() {
		return 0, errors2.WithMessage(ErrTokenInBlackList, op)
	}

	return claims.UID, nil
}

func (s *Service) Refresh(
	ctx context.Context,
	RefreshToken string,
) (string, string, error) {
	const op = "services.pcClub.auth.Refresh"

	claims, err := jwt.ParseToken(
		RefreshToken,
		s.cfg.Refresh.Secret,
	)
	if err != nil {
		return "", "", errors2.WithMessage(TokenError(err), op, "failed to parse token")
	}

	banedTokenExpirationTimeString, err := s.redisProvider.StringValue(
		ctx,
		fmt.Sprintf("%s:%d", RefreshRedisBlackListName, claims.UID),
	)
	if errors.Is(err, redis.ErrNotFound) {
		err := s.versionOwner.IncRefreshVersion(ctx, claims.UID)
		if errors.Is(err, mssql.ErrNotFound) {
			return "", "", errors2.WithMessage(TokenError(err), op, "user to update version not found")
		}
		if err != nil {
			return "", "", errors2.WithMessage(err, op, "failed to update refresh version")
		}

		access, refresh, err := generateTokens(
			claims.UID,
			claims.Version+1,
			s.cfg,
		)
		if err != nil {
			return "", "", errors2.WithMessage(err, op, "failed to generate tokens")
		}

		return access, refresh, nil
	}
	if err != nil {
		return "", "", errors2.WithMessage(err, op, "failed to find token in redis")
	}

	var banedTokenExpirationTime int64
	banedTokenExpirationTime, err = strconv.ParseInt(banedTokenExpirationTimeString, 10, 64)
	if err != nil {
		return "", "", errors2.WithMessage(err, op, "failed to parse banned token expiration time")
	}

	if banedTokenExpirationTime >= claims.ExpiresAt.Unix() {
		return "", "", errors2.WithMessage(ErrTokenInBlackList, op)
	}

	version, err := s.versionProvider.RefreshVersion(ctx, claims.UID)
	if errors.Is(err, mssql.ErrNotFound) {
		return "", "", errors2.WithMessage(ErrUserNotFound, op, "users version of refresh token not found")
	}
	if err != nil {
		return "", "", errors2.WithMessage(err, op, "failed to find refresh token version")
	}

	if version != claims.Version {
		return "", "", errors2.WithMessage(ErrInvalidRefreshVersion, op, "refresh token version is not not equal to db")
	}

	err = s.versionOwner.IncRefreshVersion(ctx, claims.UID)
	if errors.Is(err, mssql.ErrNotFound) {
		return "", "", errors2.WithMessage(ErrUserNotFound, op, "user to update version not found")
	}
	if err != nil {
		return "", "", errors2.WithMessage(err, op, "failed to update refresh version")
	}

	access, refresh, err := generateTokens(claims.UID, claims.Version+1, s.cfg)
	if err != nil {
		return "", "", errors2.WithMessage(err, op, "failed to generate tokens")
	}

	return access, refresh, nil
}

func (s *Service) Tokens(
	ctx context.Context,
	uid int64,
) (string, string, error) {
	const op = "services.pcClub.auth.Tokens"

	version, err := s.versionProvider.RefreshVersion(ctx, uid)
	if errors.Is(err, mssql.ErrNotFound) {
		return "", "", errors2.WithMessage(ErrUserNotFound, op, "users version of refresh token not found")
	}
	if err != nil {
		return "", "", errors2.WithMessage(err, op, "failed to find refresh token version")
	}

	err = s.versionOwner.IncRefreshVersion(ctx, uid)
	if errors.Is(err, mssql.ErrNotFound) {
		return "", "", errors2.WithMessage(ErrUserNotFound, op, "user to update version not found")
	}
	if err != nil {
		return "", "", errors2.WithMessage(err, op, "failed to update refresh version")
	}

	access, refresh, err := generateTokens(uid, version+1, s.cfg)
	if err != nil {
		return "", "", errors2.WithMessage(err, op, "failed to generate tokens")
	}

	return access, refresh, nil
}

func (s *Service) BanTokens(
	ctx context.Context,
	accessToken string,
	refreshToken string,
) (int64, error) {
	const op = "services.pcClub.auth.Tokens"

	claims, err := jwt.ParseToken(refreshToken, s.cfg.Refresh.Secret)
	if err != nil {
		return 0, errors2.WithMessage(TokenError(err), op, "failed to parse token")
	}

	err = s.redisOwner.SetStringWithCustomTTL(
		ctx,
		fmt.Sprintf("%s:%d", RefreshRedisBlackListName, claims.UID),
		strconv.FormatInt(claims.ExpiresAt.Unix(), 10),
		s.cfg.Refresh.TTL,
	)
	if err != nil {
		return 0, errors2.WithMessage(err, op, "failed to push refresh in black list")
	}

	claims, err = jwt.ParseToken(accessToken, s.cfg.Access.Secret)
	if err != nil {
		return 0, errors2.WithMessage(TokenError(err), op, "failed to parse token")
	}

	err = s.redisOwner.SetStringWithCustomTTL(
		ctx,
		fmt.Sprintf("%s:%d", AccessRedisBlackListName, claims.UID),
		strconv.FormatInt(claims.ExpiresAt.Unix(), 10),
		s.cfg.Access.TTL,
	)
	if err != nil {
		return 0, errors2.WithMessage(err, op, "failed to push access in black list")
	}

	return claims.UID, nil
}
