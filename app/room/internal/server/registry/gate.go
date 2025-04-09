package registry

import (
	room "github.com/go-pantheon/roma/app/room/internal/app/room/gate/registry"
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
