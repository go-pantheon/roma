package userobj

import (
	basicobj "github.com/go-pantheon/roma/app/player/internal/app/basic/gate/domain/object"
	heroobj "github.com/go-pantheon/roma/app/player/internal/app/hero/gate/domain/object"
	plunderobj "github.com/go-pantheon/roma/app/player/internal/app/plunder/gate/domain/object"
	rechargeobj "github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/domain/object"
	roomobj "github.com/go-pantheon/roma/app/player/internal/app/room/gate/domain/object"
	statusobj "github.com/go-pantheon/roma/app/player/internal/app/status/gate/domain/object"
	storageobj "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain/object"
	systemobj "github.com/go-pantheon/roma/app/player/internal/app/system/gate/domain/object"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

func (u *User) Dev() *Dev {
	return u.Modules[ModuleKey].(*Dev)
}

func (u *User) Basic() *basicobj.Basic {
	return u.Modules[basicobj.ModuleKey].(*basicobj.Basic)
}

func (u *User) Status() *statusobj.Status {
	return u.Modules[statusobj.ModuleKey].(*statusobj.Status)
}

func (u *User) HeroList() *heroobj.HeroList {
	return u.Modules[heroobj.ModuleKey].(*heroobj.HeroList)
}

func (u *User) PlunderList() *plunderobj.PlunderList {
	return u.Modules[plunderobj.ModuleKey].(*plunderobj.PlunderList)
}

func (u *User) Recharge() *rechargeobj.Recharge {
	return u.Modules[rechargeobj.ModuleKey].(*rechargeobj.Recharge)
}

func (u *User) System() *systemobj.System {
	return u.Modules[systemobj.ModuleKey].(*systemobj.System)
}

func (u *User) Storage() *storageobj.Storage {
	return u.Modules[storageobj.ModuleKey].(*storageobj.Storage)
}

func (u *User) Room() *roomobj.Room {
	return u.Modules[roomobj.ModuleKey].(*roomobj.Room)
}

func (u *User) GetModule(mod life.ModuleKey) life.Module {
	return u.Modules[mod]
}

func (u *User) RegisterModule(mod life.ModuleKey, module life.Module) {
	u.Modules[mod] = module
}
