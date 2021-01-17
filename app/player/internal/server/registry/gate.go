package registry

import (
	user "github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/registry"
)

type GateRegistrars struct {
	Rgs []Registrar
}

func NewGateRegistrars(
	user *user.UserRegistrar,
) *GateRegistrars {
	return &GateRegistrars{
		Rgs: []Registrar{
			user,
		},
	}
}
