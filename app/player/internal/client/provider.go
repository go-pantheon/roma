package client

import (
	"github.com/go-pantheon/roma/app/player/internal/client/self"
	"github.com/go-pantheon/roma/pkg/client/gate"
	"github.com/go-pantheon/roma/pkg/client/room"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewDiscovery,
	self.ProviderSet,
	gate.ProviderSet,
	room.ProviderSet,
)
