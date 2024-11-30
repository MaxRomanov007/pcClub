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
	"server/internal/http-server/middleware/auth/authAdmin"
	"server/internal/http-server/middleware/auth/authorization"
	"server/internal/http-server/middleware/logger"
	"server/internal/lib/api/logger/sl"
)

type App struct {
	Log         *slog.Logger
	HTTPSServer *http.Server
	Cfg         *config.HTTPSServerConfig
}

func New(cfg *config.HTTPSServerConfig, api *pcCLub.API) *App {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Depth", "UserService-Agent", "X-File-Size", "X-Requested-With", "If-Modified-Since", "X-File-Name", "Cache-Control", "Access-Control-Expose-Headers", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.RequestID)
	r.Use(logger.New(api.Log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Post("/register", api.Register())
	r.Post("/login", api.Login())
	r.Post("/refresh", api.Refresh())
	r.Post("/logout", api.Logout())

	r.Get("/pc-types", api.PcTypes())
	r.Get("/pcs", api.Pcs())
	r.Get("/pc-type/{type-id}", api.PcType())

	r.Get("/pc-room/{room-id}", api.PcRoom())

	r.Get("/monitor-producers", api.MonitorProducers())
	r.Get("/monitors", api.Monitors())

	r.Get("/processor-producers", api.ProcessorProducers())
	r.Get("/processors", api.Processors())

	r.Get("/video-card-producers", api.VideoCardProducers())
	r.Get("/video-cards", api.VideoCards())

	r.Get("/ram-types", api.RamTypes())
	r.Get("/ram", api.Rams())

	r.Get("/dishes", api.Dishes())
	r.Get("/dish/{dish-id}", api.Dish())

	//routes to be authorized
	r.Group(func(r chi.Router) {
		r.Use(authorization.Authorize(api.Log, api.AuthService))

		r.Post("/user", api.User())
	})

	//admin routes
	r.Group(func(r chi.Router) {
		r.Use(authAdmin.AuthAdmin(api.Log, api.AuthService, api.UserService))

		r.Post("/save-pc", api.SavePc())
		r.Post("/save-pc-type", api.SavePcType())
		r.Post("/update-pc-type", api.UpdatePcType())
		r.Post("/update-pc", api.UpdatePc())
		r.Post("/delete-pc-type", api.DeletePcType())
		r.Post("/delete-pc", api.DeletePc())

		r.Post("/save-pc-room", api.SavePcRoom())
		r.Post("/update-pc-room", api.UpdatePcRoom())
		r.Post("/delete-pc-room", api.DeletePcRoom())

		r.Post("/save-monitor-producer", api.SaveMonitorProducer())
		r.Post("/save-monitor", api.SaveMonitor())
		r.Post("/delete-monitor-producer", api.DeleteMonitorProducer())
		r.Post("/delete-monitor", api.DeleteMonitor())

		r.Post("/save-processor-producer", api.SaveProcessorProducer())
		r.Post("/save-processor", api.SaveProcessor())
		r.Post("/delete-processor-producer", api.DeleteProcessorProducer())
		r.Post("/delete-processor", api.DeleteProcessor())

		r.Post("/save-video-card-producer", api.SaveVideoCardProducer())
		r.Post("/save-video-card", api.SaveVideoCard())
		r.Post("/delete-video-card-producer", api.DeleteVideoCardProducer())
		r.Post("/delete-video-card", api.DeleteVideoCard())

		r.Post("/save-ram-type", api.SaveRamType())
		r.Post("/save-ram", api.SaveRam())
		r.Post("/delete-ram-type", api.DeleteRamType())
		r.Post("/delete-ram", api.DeleteRam())

		r.Post("/save-dish", api.SaveDish())
		r.Post("/update-dish", api.UpdateDish())
		r.Post("/delete-dish", api.DeleteDish())
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
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
