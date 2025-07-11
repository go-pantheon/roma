// Code generated by gen-api-client. DO NOT EDIT.

package service

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/google/wire"
)

var RoomServicesProviderSet = wire.NewSet(
	NewRoomServices,
)

type RoomServices struct {
	Room climsg.RoomServiceServer
}

func NewRoomServices(
	Room climsg.RoomServiceServer,
) *RoomServices {
	s := &RoomServices{}
	s.Room = Room

	return s
}
