package plunder

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/plunder/gate/domain"
)

var ProviderSet = wire.NewSet(domain.ProviderSet)
