package app

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room"
)

var ProviderSet = wire.NewSet(room.ProviderSet)
