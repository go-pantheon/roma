//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos"
	"github.com/go-kratos/kratos/log"
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/client"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/conf"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/core"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/data"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/intra"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/server"
	"github.com/vulcan-frame/vulcan-net/health"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Label, *conf.Registry, *conf.Data, log.Logger, *health.Server) (*kratos.App, func(), error) {
	panic(wire.Build(core.ProviderSet, intra.ProviderSet, data.ProviderSet, server.ProviderSet, client.ProviderSet, app.ProviderSet, newApp))
}
