package status

import (
	"github.com/go-pantheon/roma/app/player/internal/app/status/gate"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	gate.ProviderSet,
)
