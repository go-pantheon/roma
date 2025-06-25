package room

import (
	roomv1 "github.com/go-pantheon/roma/gen/api/server/room/service/room/v1"
)

func NewClient(conn *Conn) roomv1.RoomServiceClient {
	return roomv1.NewRoomServiceClient(conn.ClientConnInterface)
}
