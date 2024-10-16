package pcClub

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
	"server/internal/config"
	"server/internal/http-server/handlers/pcCLub"
	"server/internal/http-server/middleware/logger"
	"server/internal/lib/logger/sl"
)

type App struct {
	Log        *slog.Logger
	HTTPServer *http.Server
}

func New(cfg *config.HTTPServerConfig, api *pcCLub.API) *App {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Depth", "UserService-Agent", "X-File-Size", "X-Requested-With", "If-Modified-Since", "X-File-Name", "Cache-Control", "Access-Control-Expose-Headers", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.RequestID)
	router.Use(logger.New(api.Log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/register", api.Register())
	router.Post("/login", api.Login())
	router.Post("/user", api.User())
	router.Post("/refresh", api.Refresh())
	router.Post("/logout", api.Logout())

	router.Get("/pc-types", api.PcTypes())
	router.Get("/pcs", api.Pcs())
	router.Get("/pc-type/{type_id}", api.PcType())

	router.Post("/save-pc", api.SavePc())
	router.Post("/save-pc-type", api.SavePcType())

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &App{
		Log:        api.Log,
		HTTPServer: srv,
	}
}

func (a *App) MustRun() {
	const op = "app.cars.MustRun"

	log := a.Log.With(
		slog.String("operation", op),
	)

	if err := a.Run(); err != nil {
		log.Error("failed to start server", sl.Err(err))

		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.cars.Run"

	if err := a.HTTPServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: failed to start server: %w", op, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	const op = "app.cars.Run"

	err := a.HTTPServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
