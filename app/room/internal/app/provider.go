package app

import (
	"github.com/go-pantheon/roma/app/room/internal/app/room"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(room.ProviderSet)
