package service

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/gen/app/room/service"
)

var ProviderSet = wire.NewSet(service.NewRoomServices, NewTunnelService)
