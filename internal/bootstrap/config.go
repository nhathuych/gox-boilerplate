package bootstrap

import (
	"github.com/nhathuych/gox-boilerplate/internal/config"
	"go.uber.org/fx"
)

func ConfigModule(path string) fx.Option {
	return fx.Module("config",
		fx.Provide(func() (*config.Config, error) {
			return config.Load(path)
		}),
	)
}
