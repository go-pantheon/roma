package object

import (
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = "room"
)

var _ life.Module = (*Room)(nil)

type Room struct {
	Id        int64
	IsCreator bool
}

func NewRoom() *Room {
	o := &Room{}
	o.Register()
	return o
}

func (o *Room) Register() {
	userregister.Register(ModuleKey, o)
}

func NewRoomProto() *dbv1.UserRoomProto {
	p := &dbv1.UserRoomProto{}
	return p
}

func (o *Room) Marshal() ([]byte, error) {
	p := dbv1.UserRoomProtoPool.Get()
	defer dbv1.UserRoomProtoPool.Put(p)

	p.Id = o.Id
	p.IsCreator = o.IsCreator

	return proto.Marshal(p)
}

func (o *Room) Unmarshal(data []byte) error {
	p := dbv1.UserRoomProtoPool.Get()
	defer dbv1.UserRoomProtoPool.Put(p)

	if err := proto.Unmarshal(data, p); err != nil {
		return err
	}

	o.Id = p.Id
	o.IsCreator = p.IsCreator

	return nil
}

func (o *Room) EncodeClient() *climsg.UserRoomProto {
	p := &climsg.UserRoomProto{}
	p.RoomId = o.Id
	return p
}
