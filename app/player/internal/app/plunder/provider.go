package plunder

import (
	"github.com/go-pantheon/roma/app/player/internal/app/plunder/gate/domain"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(domain.ProviderSet)
