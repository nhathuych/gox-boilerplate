package bootstrap

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhathuych/gox-boilerplate/internal/config"
	apirest "github.com/nhathuych/gox-boilerplate/internal/delivery/http"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func HTTPModule() fx.Option {
	return fx.Module("http",
		fx.Provide(
			apirest.NewRouter,
			newHTTPServer,
		),
		fx.Invoke(registerHTTPServer),
	)
}

func newHTTPServer(r *gin.Engine, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
}

func registerHTTPServer(lc fx.Lifecycle, srv *http.Server, cfg *config.Config, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("http server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})
}
