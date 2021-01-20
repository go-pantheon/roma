package system

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/system/gate"
)

var ProviderSet = wire.NewSet(gate.ProviderSet)
