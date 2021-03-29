package client

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/client/gate"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/client/player"
)

var ProviderSet = wire.NewSet(
	NewDiscovery,
	gate.ProviderSet,
	player.ProviderSet,
)
