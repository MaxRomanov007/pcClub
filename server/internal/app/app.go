package app

import (
	"context"
	"fmt"
	"log/slog"
	pcClubApp "server/internal/app/pcCLub"
	"server/internal/config"
	pcClubServer "server/internal/http-server/handlers/pcCLub"
	"server/internal/services/pcClub/auth"
	"server/internal/services/pcClub/pc"
	"server/internal/services/pcClub/user"
	"server/internal/storage/redis"
	"server/internal/storage/ssms"
)

type App struct {
	PCClub *pcClubApp.App
}

func MustLoad(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	app, err := New(ctx, log, cfg)
	if err != nil {
		panic(err)
	}

	return app
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) (*App, error) {
	const op = "app.New"

	sqlserverStorage, err := ssms.New(cfg.Database.SQLServer)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create sqlserver storage: %w", op, err)
	}

	redisStorage, err := redis.New(ctx, cfg.Database.Redis)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create redis storage: %w", op, err)
	}

	authService := auth.New(cfg.Auth, redisStorage, redisStorage, sqlserverStorage, sqlserverStorage)
	userService := user.New(cfg.User, sqlserverStorage, sqlserverStorage)
	pcService := pc.New(redisStorage, redisStorage, sqlserverStorage, sqlserverStorage)

	pcClubApi := pcClubServer.New(log, cfg, userService, authService, pcService)

	pcClubApplication := pcClubApp.New(cfg.HttpServer, pcClubApi)

	return &App{
		PCClub: pcClubApplication,
	}, nil
}
