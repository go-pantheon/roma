package client

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/client/gate"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/client/room"
)

var ProviderSet = wire.NewSet(
	NewDiscovery,
	gate.ProviderSet,
	room.ProviderSet,
)
