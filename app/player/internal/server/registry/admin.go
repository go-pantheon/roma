package registry

import (
	gamedata "github.com/go-pantheon/roma/app/player/internal/app/gamedata/admin/registry"
	recharge "github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/registry"
	storage "github.com/go-pantheon/roma/app/player/internal/app/storage/admin/registry"
	user "github.com/go-pantheon/roma/app/player/internal/app/user/admin/registry"
)

type AdminRegistrars struct {
	Rgs []Registrar
}

func NewAdminRegistrars(
	user *user.UserRegistrar,
	gamedata *gamedata.GamedataRegistrar,
	storage *storage.StorageRegistrar,
	recharge *recharge.RechargeRegistrar,
) *AdminRegistrars {
	return &AdminRegistrars{
		Rgs: []Registrar{
			user,
			storage,
			gamedata,
			recharge,
		},
	}
}
