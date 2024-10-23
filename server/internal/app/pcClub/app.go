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
	Log         *slog.Logger
	HTTPSServer *http.Server
	Cfg         *config.HTTPSServerConfig
}

func New(cfg *config.HTTPSServerConfig, api *pcCLub.API) *App {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "*"},
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
	router.Get("/pc-type/{type-id}", api.PcType())

	router.Post("/save-pc", api.SavePc())
	router.Post("/save-pc-type", api.SavePcType())
	router.Post("/update-pc-type", api.UpdatePcType())
	router.Post("/update-pc", api.UpdatePc())
	router.Post("/delete-pc-type", api.DeletePcType())
	router.Post("/delete-pc", api.DeletePc())

	router.Get("/pc-room/{room-id}", api.PcRoom())
	router.Post("/save-pc-room", api.SavePcRoom())
	router.Post("/update-pc-room", api.UpdatePcRoom())
	router.Post("/delete-pc-room", api.DeletePcRoom())

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &App{
		Log:         api.Log,
		HTTPSServer: srv,
		Cfg:         cfg,
	}
}

func (a *App) MustRun() {
	const op = "app.club.MustRun"

	log := a.Log.With(
		slog.String("operation", op),
	)

	if err := a.RunClub(); err != nil {
		log.Error("failed to start server", sl.Err(err))

		panic(err)
	}
}

func (a *App) RunClub() error {
	const op = "app.club.Run"

	if err := a.HTTPSServer.ListenAndServeTLS(a.Cfg.SSLCert, a.Cfg.SSLKey); err != nil {
		return fmt.Errorf("%s: failed to start club server: %w", op, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	const op = "app.cars.Run"

	err := a.HTTPSServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to stop club server: %w", op, err)
	}

	return nil
}
