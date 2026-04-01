package bootstrap

import (
	"context"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
)

func LoggerModule() fx.Option {
	return fx.Module("logger",
		fx.Provide(func() (*zap.Logger, error) {
			// Write JSON logs to both console and file for Loki shipping via promtail.
			const logPath = "logs/app.log"
			if err := os.MkdirAll("logs", 0o755); err != nil {
				return nil, err
			}

			f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return nil, err
			}

			encoderCfg := zap.NewProductionEncoderConfig()
			encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
			encoderCfg.TimeKey = "ts"
			encoder := zapcore.NewJSONEncoder(encoderCfg)

			consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
			fileCore := zapcore.NewCore(encoder, zapcore.AddSync(f), zapcore.InfoLevel)
			return zap.New(
				zapcore.NewTee(consoleCore, fileCore),
				zap.AddCaller(),
				zap.AddStacktrace(zapcore.ErrorLevel),
			), nil
		}),
		fx.Invoke(func(lc fx.Lifecycle, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStop: func(context.Context) error {
					_ = log.Sync()
					return nil
				},
			})
		}),
	)
}
