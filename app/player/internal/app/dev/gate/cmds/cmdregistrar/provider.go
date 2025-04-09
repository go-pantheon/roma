package cmdregistrar

import (
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/storage"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/system"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/user"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewRegistrar,
	system.ProviderSet,
	storage.ProviderSet,

	user.NewAdminPlayerCommander,
	user.NewSimulateRechargeCommander,
)
