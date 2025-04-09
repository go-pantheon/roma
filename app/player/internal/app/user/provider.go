package user

import (
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	gate.ProviderSet,
	admin.ProviderSet,
)
