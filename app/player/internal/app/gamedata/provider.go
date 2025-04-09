package gamedata

import (
	"github.com/go-pantheon/roma/app/player/internal/app/gamedata/admin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(admin.ProviderSet)
