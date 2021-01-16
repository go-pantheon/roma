//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/client"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/conf"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/data"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/intra"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/server"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Label, *conf.Registry, *conf.Data, log.Logger, *health.Server) (*kratos.App, func(), error) {
	panic(wire.Build(core.ProviderSet, intra.ProviderSet, data.ProviderSet, server.ProviderSet, client.ProviderSet, app.ProviderSet, newApp))
}
