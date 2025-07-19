package app

import (
	"github.com/go-pantheon/roma/app/broadcaster/internal/app/broadcaster"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	broadcaster.ProviderSet,
)
