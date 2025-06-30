package cmdregistrar

import (
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/storage"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/system"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/user"
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
	_ *storage.ClearCommander,

	_ *user.SimulateRechargeCommander,
	_ *user.CreateAdminPlayerCommander,
) *Registrar {
	return &Registrar{}
}
