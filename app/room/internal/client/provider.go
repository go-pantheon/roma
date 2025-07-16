package client

import (
	"github.com/go-pantheon/roma/app/room/internal/client/self"
	"github.com/go-pantheon/roma/pkg/client/broadcaster"
	"github.com/go-pantheon/roma/pkg/client/player"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewDiscovery,
	self.ProviderSet,
	broadcaster.ProviderSet,
	player.ProviderSet,
)
