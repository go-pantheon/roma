package object

import (
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = "room"
)

func init() {
	userregister.Register(ModuleKey, NewRoom)
}

var _ life.Module = (*Room)(nil)

type Room struct {
	Id        int64
	IsCreator bool
}

func NewRoom() life.Module {
	o := &Room{}
	return o
}

func NewRoomProto() *dbv1.UserRoomProto {
	p := &dbv1.UserRoomProto{}
	return p
}

func (o *Room) EncodeServer() proto.Message {
	p := dbv1.UserRoomProtoPool.Get()

	p.Id = o.Id
	p.IsCreator = o.IsCreator

	return p
}

func (o *Room) DecodeServer(p proto.Message) error {
	if p == nil {
		return errors.New("room decode server nil")
	}

	op, ok := p.(*dbv1.UserRoomProto)
	if !ok {
		return errors.Errorf("room decode server invalid type: %T", p)
	}

	o.Id = op.Id
	o.IsCreator = op.IsCreator

	return nil
}

func (o *Room) EncodeClient() *climsg.UserRoomProto {
	p := &climsg.UserRoomProto{}
	p.RoomId = o.Id

	return p
}
