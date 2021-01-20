package app

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/gamedata"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/system"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user"
)

var ProviderSet = wire.NewSet(
	gamedata.ProviderSet,
	user.ProviderSet,
	system.ProviderSet,
	storage.ProviderSet,
)
