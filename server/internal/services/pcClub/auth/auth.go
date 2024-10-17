package auth

import (
	"context"
	"errors"
	"fmt"
	"server/internal/config"
	"server/internal/lib/jwt"
	"server/internal/storage/redis"
	"server/internal/storage/ssms"
	"strconv"
)

// generateTokens generates access and refresh tokens
// (1 output is access, 2 is refresh)
func generateTokens(
	uid int64,
	refreshVersion int64,
	cfg *config.AuthConfig,
) (access string, refresh string, err error) {
	access, err = jwt.NewAccessToken(uid, cfg.Access.Secret, cfg.Access.TTL)
	if err != nil {
		return "", "", err
	}

	refresh, err = jwt.NewRefreshToken(uid, refreshVersion, cfg.Refresh.Secret, cfg.Refresh.TTL)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *Service) Access(ctx context.Context, accessToken string) (int64, error) {
	const op = "services.pcClub.auth.Access"

	claims, err := jwt.ParseToken(accessToken, s.cfg.Access.Secret)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to parse token: %w", op, TokenError(err))
	}

	if !s.cfg.Access.IsBlackList {
		return claims.UID, nil
	}

	banedTokenExpirationTimeString, err := s.redisProvider.StringValue(ctx, fmt.Sprintf("%s:%d", s.cfg.Access.RedisBlackListName, claims.UID))
	if errors.Is(err, redis.ErrNotFound) {
		return claims.UID, nil
	}
	if err != nil {
		return 0, fmt.Errorf("%s: failed to find token in redis: %w", op, err)
	}

	var banedTokenExpirationTime int64
	banedTokenExpirationTime, err = strconv.ParseInt(banedTokenExpirationTimeString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to parse banned token expiration time: %w", op, err)
	}

	if banedTokenExpirationTime >= claims.ExpiresAt.Unix() {
		return 0, fmt.Errorf("%s: %w", op, ErrTokenInBlackList)
	}

	return claims.UID, nil
}

func (s *Service) Refresh(
	ctx context.Context,
	RefreshToken string,
) (string, string, error) {
	const op = "services.pcClub.auth.Refresh"

	claims, err := jwt.ParseToken(RefreshToken, s.cfg.Refresh.Secret)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to parse token: %w", op, TokenError(err))
	}

	banedTokenExpirationTimeString, err := s.redisProvider.StringValue(ctx, fmt.Sprintf("%s:%d", s.cfg.Refresh.RedisBlackListName, claims.UID))
	if errors.Is(err, redis.ErrNotFound) {
		err := s.versionOwner.UpdateRefreshVersion(ctx, claims.UID, claims.Version+1)
		if errors.Is(err, ssms.ErrNotFound) {
			return "", "", fmt.Errorf("%s: user to update version not found: %w", op, ErrUserNotFound)
		}
		if err != nil {
			return "", "", fmt.Errorf("%s: failed to update refresh version: %w", op, err)
		}

		access, refresh, err := generateTokens(claims.UID, claims.Version+1, s.cfg)
		if err != nil {
			return "", "", fmt.Errorf("%s: failed to generate tokens: %w", op, err)
		}

		return access, refresh, nil
	}
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to find token in redis: %w", op, err)
	}

	var banedTokenExpirationTime int64
	banedTokenExpirationTime, err = strconv.ParseInt(banedTokenExpirationTimeString, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to parse banned token expiration time: %w", op, err)
	}

	if banedTokenExpirationTime >= claims.ExpiresAt.Unix() {
		return "", "", fmt.Errorf("%s: %w", op, ErrTokenInBlackList)
	}

	version, err := s.versionProvider.RefreshVersion(ctx, claims.UID)
	if errors.Is(err, ssms.ErrNotFound) {
		return "", "", fmt.Errorf("%s: users version of refresh token not found: %w", op, ErrUserNotFound)
	}
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to find refresh token version: %w", op, err)
	}

	if version != claims.Version {
		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidRefreshVersion)
	}

	err = s.versionOwner.UpdateRefreshVersion(ctx, claims.UID, claims.Version+1)
	if errors.Is(err, ssms.ErrNotFound) {
		return "", "", fmt.Errorf("%s: user to update version not found: %w", op, ErrUserNotFound)
	}
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to update refresh version: %w", op, err)
	}

	access, refresh, err := generateTokens(claims.UID, claims.Version+1, s.cfg)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to generate tokens: %w", op, err)
	}

	return access, refresh, nil
}

func (s *Service) Tokens(
	ctx context.Context,
	uid int64,
) (string, string, error) {
	const op = "services.pcClub.auth.Tokens"

	version, err := s.versionProvider.RefreshVersion(ctx, uid)
	if errors.Is(err, ssms.ErrNotFound) {
		return "", "", fmt.Errorf("%s: users version of refresh token not found: %w", op, ErrUserNotFound)
	}
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to find refresh token version: %w", op, err)
	}

	err = s.versionOwner.UpdateRefreshVersion(ctx, uid, version+1)
	if errors.Is(err, ssms.ErrNotFound) {
		return "", "", fmt.Errorf("%s: user to update version not found: %w", op, ErrUserNotFound)
	}
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to update refresh version: %w", op, err)
	}

	access, refresh, err := generateTokens(uid, version+1, s.cfg)
	if err != nil {
		return "", "", fmt.Errorf("%s: failed to generate tokens: %w", op, err)
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
		return 0, fmt.Errorf("%s: failed to parse token: %w", op, TokenError(err))
	}

	err = s.redisOwner.SetStringWithCustomTTL(
		ctx,
		fmt.Sprintf("%s:%d", s.cfg.Refresh.RedisBlackListName, claims.UID),
		strconv.FormatInt(claims.ExpiresAt.Unix(), 10),
		s.cfg.Refresh.TTL,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to push refresh in black list: %w", op, err)
	}

	if !s.cfg.Access.IsBlackList {
		return claims.UID, nil
	}

	claims, err = jwt.ParseToken(accessToken, s.cfg.Access.Secret)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to parse token: %w", op, TokenError(err))
	}

	err = s.redisOwner.SetStringWithCustomTTL(
		ctx,
		fmt.Sprintf("%s:%d", s.cfg.Access.RedisBlackListName, claims.UID),
		strconv.FormatInt(claims.ExpiresAt.Unix(), 10),
		s.cfg.Access.TTL,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to push access in black list: %w", op, err)
	}

	return claims.UID, nil
}
