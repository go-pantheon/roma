package room

import (
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(gate.ProviderSet, admin.ProviderSet)
