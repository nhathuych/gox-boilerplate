package bootstrap

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/nhathuych/gox-boilerplate/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func RabbitMQModule() fx.Option {
	return fx.Module("rabbitmq",
		fx.Provide(NewAMQPConnection),
		fx.Invoke(registerAMQPLifecycle),
	)
}

func NewAMQPConnection(cfg *config.Config, log *zap.Logger) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, err
	}
	log.Info("rabbitmq connected")
	return conn, nil
}

func registerAMQPLifecycle(lc fx.Lifecycle, conn *amqp.Connection, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("closing rabbitmq connection")
			return conn.Close()
		},
	})
}
