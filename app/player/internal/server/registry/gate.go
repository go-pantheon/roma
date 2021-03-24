package registry

import (
	dev "github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/registry"
	hero "github.com/vulcan-frame/vulcan-game/app/player/internal/app/hero/gate/registry"
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
	dev *dev.DevRegistrar,
	storage *storage.StorageRegistrar,
	hero *hero.HeroRegistrar,
) *GateRegistrars {
	return &GateRegistrars{
		Rgs: []Registrar{
			system,
			user,
			dev,
			storage,
			hero,
		},
	}
}
