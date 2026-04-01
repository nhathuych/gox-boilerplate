package bootstrap

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/nhathuych/gox-boilerplate/internal/config"
	"github.com/nhathuych/gox-boilerplate/internal/worker"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func WorkerProvidersModule() fx.Option {
	return fx.Module("worker",
		fx.Provide(func(cfg *config.Config, rdb *redis.Client, log *zap.Logger) *worker.MailWorker {
			return worker.NewMailWorker(cfg, rdb, log)
		}),
		fx.Invoke(func(lc fx.Lifecycle, log *zap.Logger, conn *amqp.Connection, w *worker.MailWorker) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					_ = ctx
					go func() {
						if err := w.Run(conn); err != nil {
							log.Error("mail worker stopped", zap.Error(err))
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return w.Shutdown(ctx)
				},
			})
		}),
	)
}
