package gamedata

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/gamedata/admin"
)

var ProviderSet = wire.NewSet(admin.ProviderSet)
