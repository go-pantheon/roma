package client

import (
	"github.com/go-pantheon/roma/app/room/internal/client/gate"
	"github.com/go-pantheon/roma/app/room/internal/client/player"
	"github.com/go-pantheon/roma/app/room/internal/client/self"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewDiscovery,
	self.ProviderSet,
	gate.ProviderSet,
	player.ProviderSet,
)
