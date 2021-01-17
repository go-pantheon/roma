package service

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/gen/app/player/service"
)

var ProviderSet = wire.NewSet(service.NewPlayerServices, NewTunnelService)
