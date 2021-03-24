package cmdregistrar

import (
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds/storage"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds/system"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds/user"
)

type Registrar struct {
}

// NewRegistrar add all commands as parameters, and register them to wire to call NewXXCommander()
func NewRegistrar(
	_ *system.ShowTimeCommander,
	_ *system.ChangeTimeCommander,

	_ *storage.AddItemCommander,
	_ *storage.SubItemCommander,
	_ *storage.AddPackCommander,
	_ *storage.SubPackCommander,

	_ *user.SimulateRechargeCommander,
	_ *user.CreateAdminPlayerCommander,
) *Registrar {
	return &Registrar{}
}
