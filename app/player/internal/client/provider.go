package client

import (
	"github.com/go-pantheon/roma/app/player/internal/client/gate"
	"github.com/go-pantheon/roma/app/player/internal/client/room"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewDiscovery,
	gate.ProviderSet,
	room.ProviderSet,
)
