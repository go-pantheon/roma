package core

import (
	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/go-pantheon/fabrica-util/errors"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/room/intra/v1"
	"github.com/go-pantheon/roma/gen/app/room/handler"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

func NewResponser(sendFunc life.SendFunc) life.Responsive {
	return life.NewResponser(sendFunc,
		func(mod int32, seq int32, obj int64, sc proto.Message) (xnet.TunnelMessage, error) {
			return handler.TakeProtoRoomTunnelResponse(0, mod, seq, obj, sc)
		},
	)
}

func SendFunc(stream intrav1.TunnelService_TunnelServer, p xnet.TunnelMessage) error {
	msg, ok := p.(*intrav1.TunnelResponse)
	if !ok {
		return errors.New("intrav1.TunnelResponse proto type conversion failed")
	}

	defer handler.PutRoomTunnelResponse(msg)

	if err := stream.Send(msg); err != nil {
		return errors.Wrapf(err, "intrav1.TunnelResponse send failed")
	}

	return nil
}

func mockSendFunc(out xnet.TunnelMessage) error {
	handler.PutRoomTunnelResponse(out.(*intrav1.TunnelResponse))
	return nil
}
