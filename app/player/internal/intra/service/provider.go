package service

import (
	"github.com/go-pantheon/roma/gen/app/player/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(service.NewPlayerServices, NewTunnelService)
