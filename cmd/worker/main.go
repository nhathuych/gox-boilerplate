package main

import (
	"flag"
	"os"

	"github.com/nhathuych/gox-boilerplate/internal/bootstrap"
	"go.uber.org/fx"
)

func main() {
	cfgPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()
	path := *cfgPath
	if path == "" {
		path = os.Getenv("GOX_CONFIG")
	}
	fx.New(
		bootstrap.WorkerModule(path),
		fx.NopLogger,
	).Run()
}
