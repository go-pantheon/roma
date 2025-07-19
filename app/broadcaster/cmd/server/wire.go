//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-net/http/health"
	"github.com/go-pantheon/roma/app/broadcaster/internal/app"
	"github.com/go-pantheon/roma/app/broadcaster/internal/conf"
	"github.com/go-pantheon/roma/app/broadcaster/internal/core"
	"github.com/go-pantheon/roma/app/broadcaster/internal/data"
	"github.com/go-pantheon/roma/app/broadcaster/internal/server"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(srv *conf.Server, label *conf.Label, registry *conf.Registry, dataConf *conf.Data, logger log.Logger, health *health.Server) (*kratos.App, func(), error) {
	panic(wire.Build(core.ProviderSet, data.ProviderSet, server.ProviderSet, app.ProviderSet, newApp))
}
