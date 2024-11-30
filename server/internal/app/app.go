package app

import (
	"context"
	"fmt"
	"log/slog"
	pcClubApp "server/internal/app/pcCLub"
	"server/internal/config"
	pcClubServer "server/internal/http-server/handlers/pcCLub"
	"server/internal/services/pcClub/auth"
	"server/internal/services/pcClub/components/monitor"
	"server/internal/services/pcClub/components/processor"
	"server/internal/services/pcClub/components/ram"
	"server/internal/services/pcClub/components/videoCard"
	"server/internal/services/pcClub/dish"
	"server/internal/services/pcClub/pc"
	"server/internal/services/pcClub/pcRoom"
	"server/internal/services/pcClub/pcType"
	"server/internal/services/pcClub/user"
	gorm "server/internal/storage/mssql"
	"server/internal/storage/redis"
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

	mssqlStorage, err := gorm.New(cfg.Database.SQLServer)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create mssql storage: %w", op, err)
	}

	redisStorage, err := redis.New(ctx, cfg.Database.Redis)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create redis storage: %w", op, err)
	}

	authService := auth.New(cfg.Auth, redisStorage, redisStorage, mssqlStorage, mssqlStorage)
	userService := user.New(cfg.User, mssqlStorage, mssqlStorage)
	pcTypeService := pcType.New(mssqlStorage, mssqlStorage, redisStorage, redisStorage)
	pcService := pc.New(mssqlStorage, mssqlStorage)
	pcRoomService := pcRoom.New(redisStorage, redisStorage, mssqlStorage, mssqlStorage)
	processorService := processor.New(mssqlStorage, mssqlStorage, redisStorage, redisStorage)
	monitorService := monitor.New(mssqlStorage, mssqlStorage, redisStorage, redisStorage)
	videoCardService := videoCard.New(mssqlStorage, mssqlStorage, redisStorage, redisStorage)
	ramService := ram.New(mssqlStorage, mssqlStorage, redisStorage, redisStorage)
	dishService := dish.New(mssqlStorage, mssqlStorage, redisStorage, redisStorage)

	pcClubApi := pcClubServer.New(
		log,
		cfg,
		userService,
		authService,
		pcTypeService,
		pcService,
		pcRoomService,
		pcClubServer.ComponentsService{
			Processor: processorService,
			Monitor:   monitorService,
			VideoCard: videoCardService,
			Ram:       ramService,
		},
		dishService,
	)

	pcClubApplication := pcClubApp.New(cfg.HttpsServer, pcClubApi)

	return &App{
		PCClub: pcClubApplication,
	}, nil
}
