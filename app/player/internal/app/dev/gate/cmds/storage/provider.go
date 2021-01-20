package storage

import "github.com/google/wire"

const Mod = "Storage"

var ProviderSet = wire.NewSet(
	NewAddItemCommander,
	NewSubItemCommander,
	NewAddPackCommander,
	NewSubPackCommander,
)
