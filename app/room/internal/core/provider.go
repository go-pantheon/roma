package core

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/data"
)

var ProviderSet = wire.NewSet(NewManager, data.ProviderSet)
