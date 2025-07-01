package guildobj

import (
	"time"

	"github.com/go-pantheon/fabrica-util/xtime"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
)

type Room struct {
	Id        int64
	CreatedAt time.Time
	Version   int64
}

func NewRoom() *Room {
	u := &Room{
		CreatedAt: time.Now(),
	}

	return u
}

func (o *Room) EncodeServer() *dbv1.RoomProto {
	p := &dbv1.RoomProto{
		Id:        o.Id,
		CreatedAt: o.CreatedAt.Unix(),
		Version:   o.Version,
	}

	return p
}

func (o *Room) DecodeServer(p *dbv1.RoomProto) (err error) {
	if p == nil {
		return
	}

	o.Id = p.Id
	o.CreatedAt = xtime.Time(p.CreatedAt)
	o.Version = p.Version

	return nil
}
