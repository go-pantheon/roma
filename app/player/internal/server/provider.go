package server

import (
	"github.com/go-pantheon/roma/app/player/internal/server/registry"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(registry.ProviderSet, NewGRPCServer, NewHTTPServer, NewRegistrar)
