package system

import (
	"github.com/go-pantheon/roma/app/player/internal/app/system/gate"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(gate.ProviderSet)
