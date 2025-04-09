package server

import (
	"github.com/go-pantheon/roma/app/room/internal/server/registry"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	registry.ProviderSet,
	NewGRPCServer,
	NewHTTPServer,
	NewRegistrar,
)
