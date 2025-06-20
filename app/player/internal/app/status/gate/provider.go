package gate

import (
	"github.com/go-pantheon/roma/app/player/internal/app/status/gate/domain"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	domain.ProviderSet,
)
