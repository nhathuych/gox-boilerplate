package bootstrap

import (
	"context"
	"time"

	"github.com/nhathuych/gox-boilerplate/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func RedisModule() fx.Option {
	return fx.Module("redis",
		fx.Provide(NewRedis),
		fx.Invoke(registerRedisPing),
		fx.Invoke(registerRedisLifecycle),
	)
}

func registerRedisPing(lc fx.Lifecycle, rdb *redis.Client, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			if err := rdb.Ping(pingCtx).Err(); err != nil {
				return err
			}
			log.Info("redis ping ok")
			return nil
		},
	})
}

func NewRedis(cfg *config.Config, log *zap.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	})
	log.Info("redis client configured", zap.String("addr", cfg.Redis.Addr))
	return rdb
}

func registerRedisLifecycle(lc fx.Lifecycle, rdb *redis.Client, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("closing redis")
			return rdb.Close()
		},
	})
}
