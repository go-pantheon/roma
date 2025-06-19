package userobj

import (
	basicobj "github.com/go-pantheon/roma/app/player/internal/app/basic/gate/domain/object"
	devobj "github.com/go-pantheon/roma/app/player/internal/app/dev/gate/domain/object"
	heroobj "github.com/go-pantheon/roma/app/player/internal/app/hero/gate/domain/object"
	plunderobj "github.com/go-pantheon/roma/app/player/internal/app/plunder/gate/domain/object"
	rechargeobj "github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/domain/object"
	roomobj "github.com/go-pantheon/roma/app/player/internal/app/room/gate/domain/object"
	statusobj "github.com/go-pantheon/roma/app/player/internal/app/status/gate/domain/object"
	storageobj "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain/object"
	systemobj "github.com/go-pantheon/roma/app/player/internal/app/system/gate/domain/object"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

func (u *User) Dev() *devobj.Dev {
	return u.modules[devobj.ModuleKey].(*devobj.Dev)
}

func (u *User) Basic() *basicobj.Basic {
	return u.modules[basicobj.ModuleKey].(*basicobj.Basic)
}

func (u *User) Status() *statusobj.Status {
	return u.modules[statusobj.ModuleKey].(*statusobj.Status)
}

func (u *User) HeroList() *heroobj.HeroList {
	return u.modules[heroobj.ModuleKey].(*heroobj.HeroList)
}

func (u *User) PlunderList() *plunderobj.PlunderList {
	return u.modules[plunderobj.ModuleKey].(*plunderobj.PlunderList)
}

func (u *User) Recharge() *rechargeobj.Recharge {
	return u.modules[rechargeobj.ModuleKey].(*rechargeobj.Recharge)
}

func (u *User) System() *systemobj.System {
	return u.modules[systemobj.ModuleKey].(*systemobj.System)
}

func (u *User) Storage() *storageobj.Storage {
	return u.modules[storageobj.ModuleKey].(*storageobj.Storage)
}

func (u *User) Room() *roomobj.Room {
	return u.modules[roomobj.ModuleKey].(*roomobj.Room)
}

func (u *User) Module(mod life.ModuleKey) life.Module {
	return u.modules[mod]
}
