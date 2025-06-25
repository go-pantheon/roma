package object

import (
	"time"

	"github.com/go-pantheon/fabrica-util/xtime"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

type Room struct {
	Id        int64
	Sid       uint64
	RoomType  uint64
	Members   []*Member
	CreatedAt time.Time
	Version   int64
}

func NewRoom() *Room {
	u := &Room{
		CreatedAt: time.Now(),
	}
	return u
}

func (o *Room) EncodeServer(p *dbv1.RoomProto, keys []life.ModuleKey) {
	p.Id = o.Id
	p.Sid = o.Sid
	p.RoomType = o.RoomType
	p.Members = make([]*dbv1.RoomMemberProto, 0, len(o.Members))
	p.CreatedAt = o.CreatedAt.Unix()

	for _, m := range o.Members {
		p.Members = append(p.Members, m.encodeServer())
	}
}

func (o *Room) DecodeServer(p *dbv1.RoomProto) (err error) {
	if p == nil {
		return
	}

	o.Id = p.Id
	o.Sid = p.Sid
	o.RoomType = p.RoomType
	o.CreatedAt = xtime.Time(p.CreatedAt)

	o.Members = make([]*Member, 0, len(p.Members))

	for _, mp := range p.Members {
		m := NewMember()
		o.Members = append(o.Members, m.decodeServer(mp))
	}

	return nil
}

func (o *Room) EncodeClient() *climsg.RoomProto {
	p := &climsg.RoomProto{
		Basic: &climsg.RoomBasicProto{
			Id: o.Id,
		},
		Members: make(map[int64]*climsg.RoomMemberProto),
	}

	for _, m := range o.Members {
		p.Members[m.Id] = m.EncodeClient()
	}

	return p
}
