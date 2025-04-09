package gate

import (
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/data"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/domain"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(data.ProviderSet, domain.ProviderSet)
