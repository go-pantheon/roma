package recharge

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/admin"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/gate"
)

var ProviderSet = wire.NewSet(gate.ProviderSet, admin.ProviderSet)
