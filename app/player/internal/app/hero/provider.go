package hero

import (
	"github.com/go-pantheon/roma/app/player/internal/app/hero/gate"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(gate.ProviderSet)
