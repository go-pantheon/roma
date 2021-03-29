package registry

import (
	room "github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate/registry"
)

type GateRegistrars struct {
	Rgs []Registrar
}

func NewGateRegistrars(
	room *room.RoomRegistrar,
) *GateRegistrars {
	return &GateRegistrars{
		Rgs: []Registrar{
			room,
		},
	}
}
