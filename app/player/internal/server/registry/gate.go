package registry

import (
	dev "github.com/go-pantheon/roma/app/player/internal/app/dev/gate/registry"
	hero "github.com/go-pantheon/roma/app/player/internal/app/hero/gate/registry"
	storage "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/registry"
	system "github.com/go-pantheon/roma/app/player/internal/app/system/gate/registry"
	user "github.com/go-pantheon/roma/app/player/internal/app/user/gate/registry"
)

type GateRegistrars struct {
	Rgs []Registrar
}

func NewGateRegistrars(
	user *user.UserRegistrar,
	dev *dev.DevRegistrar,
	system *system.SystemRegistrar,
	storage *storage.StorageRegistrar,
	hero *hero.HeroRegistrar,
) *GateRegistrars {
	return &GateRegistrars{
		Rgs: []Registrar{
			user,
			dev,
			system,
			storage,
			hero,
		},
	}
}
