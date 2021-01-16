package app

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user"
)

var ProviderSet = wire.NewSet(
	user.ProviderSet,
)
