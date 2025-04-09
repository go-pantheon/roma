package service

import (
	"github.com/go-pantheon/roma/gen/app/room/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(service.NewRoomServices, NewTunnelService)
