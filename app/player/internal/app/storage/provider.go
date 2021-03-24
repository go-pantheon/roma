package storage

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/admin"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate"
)

var ProviderSet = wire.NewSet(
	gate.ProviderSet,
	admin.ProviderSet,
)
