package object

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

type Room struct {
	Id        int64
	IsCreator bool
}

func NewRoom() *Room {
	return &Room{}
}

func NewRoomProto() *dbv1.UserRoomProto {
	p := &dbv1.UserRoomProto{}
	return p
}

func (o *Room) EncodeServer(p *dbv1.UserRoomProto) {
	p.Id = o.Id
	p.IsCreator = o.IsCreator
}

func (o *Room) DecodeServer(p *dbv1.UserRoomProto) {
	if p == nil {
		return
	}
	o.Id = p.Id
	o.IsCreator = p.IsCreator
}

func (o *Room) EncodeClient() *climsg.UserRoomProto {
	p := &climsg.UserRoomProto{}
	p.RoomId = o.Id
	return p
}
