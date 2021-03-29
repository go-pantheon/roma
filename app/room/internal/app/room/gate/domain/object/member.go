package object

import (
	"time"

	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/room/v1"
)

type Member struct {
	Id        int64
	JoinedAt  time.Time
	CreatedAt time.Time
	Version   int64
}

func NewMember() *Member {
	return &Member{
		CreatedAt: time.Now(),
	}
}

func NewMemberProto() *dbv1.RoomMemberProto {
	return &dbv1.RoomMemberProto{}
}

func (o *Member) EncodeServer(p *dbv1.RoomMemberProto) {
	p.Id = o.Id
	p.JoinedAt = o.JoinedAt.Unix()
}

func (o *Member) DecodeServer(p *dbv1.RoomMemberProto) {
	o.Id = p.Id
	o.JoinedAt = time.Unix(p.JoinedAt, 0)
}

func (o *Member) EncodeClient() *climsg.RoomMemberProto {
	return &climsg.RoomMemberProto{
		JoinedAt: o.JoinedAt.Unix(),
	}
}
