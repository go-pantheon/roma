package registry

import (
	storage "github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/registry"
	system "github.com/vulcan-frame/vulcan-game/app/player/internal/app/system/gate/registry"
	user "github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/registry"
)

type GateRegistrars struct {
	Rgs []Registrar
}

func NewGateRegistrars(
	system *system.SystemRegistrar,
	user *user.UserRegistrar,
	storage *storage.StorageRegistrar,
) *GateRegistrars {
	return &GateRegistrars{
		Rgs: []Registrar{
			system,
			user,
			storage,
		},
	}
}
