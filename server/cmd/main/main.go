package main

import (
	"context"
	"github.com/fatih/color"
	"log/slog"
	"os"
	"os/signal"
	"server/internal/app"
	"server/internal/config"
	"server/internal/lib/api/logger/handlers/slogpretty"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	application := app.MustLoad(ctx, log, cfg)

	log.Info("starting server...")

	go application.PCClub.MustRun()

	log.Info("server running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping pcClub server", slog.String("signal", sign.String()))

	if err := application.PCClub.Stop(ctx); err != nil {
		log.Error("failed stop pcClub server")

		panic(err)
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	color.NoColor = false

	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
