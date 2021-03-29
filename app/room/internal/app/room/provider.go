package room

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/admin"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate"
)

var ProviderSet = wire.NewSet(gate.ProviderSet, admin.ProviderSet)
