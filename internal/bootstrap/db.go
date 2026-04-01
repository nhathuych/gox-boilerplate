package bootstrap

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhathuych/gox-boilerplate/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func DatabaseModule() fx.Option {
	return fx.Module("database",
		fx.Provide(NewPool),
		fx.Invoke(registerPoolLifecycle),
	)
}

func NewPool(cfg *config.Config, log *zap.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.ConnTimeout)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(cfg.Database.DSN)
	if err != nil {
		return nil, err
	}
	poolCfg.MaxConns = cfg.Database.MaxConns
	poolCfg.MinConns = cfg.Database.MinConns
	poolCfg.MaxConnLifetime = cfg.Database.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}
	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}
	log.Info("database connected")
	return pool, nil
}

func registerPoolLifecycle(lc fx.Lifecycle, pool *pgxpool.Pool, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("closing database pool")
			pool.Close()
			return nil
		},
	})
}
