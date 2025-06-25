package object

import (
	"time"

	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
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

func (o *Member) encodeServer() *dbv1.RoomMemberProto {
	p := dbv1.RoomMemberProtoPool.Get()

	p.Id = o.Id
	p.JoinedAt = o.JoinedAt.Unix()

	return p
}

func (o *Member) decodeServer(p *dbv1.RoomMemberProto) *Member {
	o.Id = p.Id
	o.JoinedAt = time.Unix(p.JoinedAt, 0)

	return o
}

func (o *Member) EncodeClient() *climsg.RoomMemberProto {
	return &climsg.RoomMemberProto{
		JoinedAt: o.JoinedAt.Unix(),
	}
}
