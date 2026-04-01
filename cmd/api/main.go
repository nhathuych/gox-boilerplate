package main

import (
	"flag"
	"os"

	_ "github.com/nhathuych/gox-boilerplate/docs"
	"github.com/nhathuych/gox-boilerplate/internal/bootstrap"
	"go.uber.org/fx"
)

// @title Gox Boilerplate API
// @version 1.0
// @description REST API with JWT, RBAC, and articles CRUD.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfgPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()
	path := *cfgPath
	if path == "" {
		path = os.Getenv("GOX_CONFIG")
	}
	fx.New(
		bootstrap.APIModule(path),
		fx.NopLogger,
	).Run()
}
