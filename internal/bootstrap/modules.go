package bootstrap

import "go.uber.org/fx"

// APIModule wires config, logging, PostgreSQL, Redis, domain providers, and the HTTP server.
func APIModule(configPath string) fx.Option {
	return fx.Options(
		ConfigModule(configPath),
		LoggerModule(),
		DatabaseModule(),
		RedisModule(),
		ProvidersModule(),
		HTTPModule(),
	)
}

// WorkerModule wires config, logging, Redis, RabbitMQ, and worker consumers.
func WorkerModule(configPath string) fx.Option {
	return fx.Options(
		ConfigModule(configPath),
		LoggerModule(),
		RedisModule(),
		RabbitMQModule(),
		WorkerProvidersModule(),
	)
}
