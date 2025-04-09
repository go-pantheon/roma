package client

import (
	"github.com/go-pantheon/roma/app/room/internal/client/gate"
	"github.com/go-pantheon/roma/app/room/internal/client/player"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewDiscovery,
	gate.ProviderSet,
	player.ProviderSet,
)
