package server

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/server/registry"
)

var ProviderSet = wire.NewSet(
	registry.ProviderSet,
	NewGRPCServer,
	NewHTTPServer,
	NewRegistrar,
)
