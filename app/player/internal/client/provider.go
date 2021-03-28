package client

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/client/gate"
)

var ProviderSet = wire.NewSet(
	NewDiscovery,
	gate.ProviderSet,
)
