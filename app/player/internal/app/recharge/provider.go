package recharge

import (
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(gate.ProviderSet, admin.ProviderSet)
