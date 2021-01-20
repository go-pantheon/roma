package gate

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/gate/data"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/gate/domain"
)

var ProviderSet = wire.NewSet(data.ProviderSet, domain.ProviderSet)
