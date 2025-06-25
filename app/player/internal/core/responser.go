package core

import (
	"github.com/go-pantheon/roma/gen/app/player/handler"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

func (m *Manager) NewResponser(replyFunc life.ReplyFunc) life.Responsive {
	return life.NewResponser(replyFunc,
		func(mod int32, seq int32, obj int64, sc proto.Message) (proto.Message, error) {
			return handler.NewPlayerResponseProto(mod, seq, obj, sc)
		}, func(mod int32, seq int32, obj int64, body []byte) proto.Message {
			return handler.NewPlayerResponseProtoByData(mod, seq, obj, body)
		},
	)
}

func (m *Manager) NewBroadcaster() life.Broadcastable {
	return life.NewBroadcaster(m.Pusher())
}

var mockResponseFunc = func(out proto.Message) error {
	return nil
}
