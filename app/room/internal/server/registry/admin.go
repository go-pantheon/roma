package registry

import (
	room "github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/admin/registry"
)

type AdminRegistrars struct {
	Rgs []Registrar
}

func NewAdminRegistrars(
	room *room.RoomRegistrar,
) *AdminRegistrars {
	return &AdminRegistrars{
		Rgs: []Registrar{
			room,
		},
	}
}
