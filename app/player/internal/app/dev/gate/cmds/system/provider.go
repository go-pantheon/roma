package system

import (
	"github.com/google/wire"
)

const Mod = "System"

var ProviderSet = wire.NewSet(
	NewShowTimeCommander,
	NewChangeTimeCommander,
)
