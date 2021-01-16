package user

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate"
)

var ProviderSet = wire.NewSet(
	gate.ProviderSet,
)
