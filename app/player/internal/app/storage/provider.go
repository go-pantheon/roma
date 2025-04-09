package storage

import (
	"github.com/go-pantheon/roma/app/player/internal/app/storage/admin"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	gate.ProviderSet,
	admin.ProviderSet,
)
