package registry

import (
	room "github.com/go-pantheon/roma/app/room/internal/app/room/admin/registry"
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
