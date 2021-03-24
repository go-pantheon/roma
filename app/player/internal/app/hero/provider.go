package hero

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/hero/gate"
)

var ProviderSet = wire.NewSet(gate.ProviderSet)
