package core

import (
	"github.com/go-pantheon/fabrica-net/xnet"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/player/intra/v1"
	"github.com/go-pantheon/roma/gen/app/player/handler"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func NewResponser(replyFunc life.ReplyFunc) life.Responsive {
	return life.NewResponser(replyFunc,
		func(mod int32, seq int32, obj int64, sc proto.Message) (xnet.TunnelMessage, error) {
			return handler.TakeProtoPlayerTunnelResponse(mod, seq, obj, sc)
		},
	)
}

func ReplyFunc(stream intrav1.TunnelService_TunnelServer, p xnet.TunnelMessage) error {
	msg, ok := p.(*intrav1.TunnelResponse)
	if !ok {
		return errors.New("intrav1.TunnelResponse proto type conversion failed")
	}

	defer handler.PutPlayerTunnelResponse(msg)

	if err := stream.Send(msg); err != nil {
		return errors.Wrapf(err, "intrav1.TunnelResponse send failed")
	}

	return nil
}

func mockResponseFunc(out xnet.TunnelMessage) error {
	handler.PutPlayerTunnelResponse(out.(*intrav1.TunnelResponse))
	return nil
}
