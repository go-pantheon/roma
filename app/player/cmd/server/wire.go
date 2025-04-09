//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-net/health"
	"github.com/go-pantheon/roma/app/player/internal/app"
	"github.com/go-pantheon/roma/app/player/internal/client"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/app/player/internal/data"
	"github.com/go-pantheon/roma/app/player/internal/intra"
	"github.com/go-pantheon/roma/app/player/internal/server"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Label, *conf.Recharge, *conf.Registry, *conf.Data, log.Logger, *health.Server) (*kratos.App, func(), error) {
	panic(wire.Build(core.ProviderSet, intra.ProviderSet, data.ProviderSet, server.ProviderSet, client.ProviderSet, app.ProviderSet, newApp))
}
