package cmdregistrar

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds/storage"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds/system"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds/user"
)

var ProviderSet = wire.NewSet(
	NewRegistrar,
	system.ProviderSet,
	storage.ProviderSet,

	user.NewAdminPlayerCommander,
	user.NewSimulateRechargeCommander,
)
