package life

import (
	"sync/atomic"

	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"google.golang.org/protobuf/proto"
)

type ReplyFunc = func(msg xnet.TunnelMessage) error

type Responsive interface {
	Reply(mod int32, seq int32, obj int64, sc proto.Message) error
	ReplyImmediately(mod int32, seq int32, obj int64, sc proto.Message) error
	ReplyTunnelMessage(msg xnet.TunnelMessage)
	ConsumeTunnelResponse() <-chan xnet.TunnelMessage
	ExecuteReply(msg xnet.TunnelMessage) error
	ReplyFunc() ReplyFunc
	UpdateReplyFunc(replyFunc ReplyFunc)
}

type BuildTunnelResponseFunc = func(mod int32, seq int32, obj int64, msg proto.Message) (xnet.TunnelMessage, error)

var _ Responsive = (*Responser)(nil)

type Responser struct {
	msgs      chan xnet.TunnelMessage
	replyFunc atomic.Value

	buildTunnelResponseFunc BuildTunnelResponseFunc
}

func NewResponser(replyFunc ReplyFunc, buildRespFunc BuildTunnelResponseFunc) *Responser {
	o := &Responser{
		msgs:                 make(chan xnet.TunnelMessage, constants.WorkerReplySize),
		replyFunc:            atomic.Value{},
		buildTunnelResponseFunc: buildRespFunc,
	}

	o.replyFunc.Store(replyFunc)

	return o
}

func (w *Responser) Reply(mod int32, seq int32, obj int64, sc proto.Message) error {
	msg, err := w.buildTunnelResponseFunc(mod, seq, obj, sc)
	if err != nil {
		return err
	}

	w.msgs <- msg

	return nil
}

func (w *Responser) ReplyTunnelMessage(msg xnet.TunnelMessage) {
	w.msgs <- msg
}

func (w *Responser) ReplyImmediately(mod int32, seq int32, obj int64, sc proto.Message) error {
	msg, err := w.buildTunnelResponseFunc(mod, seq, obj, sc)
	if err != nil {
		return err
	}

	return w.replyFunc.Load().(ReplyFunc)(msg)
}

func (w *Responser) ExecuteReply(out xnet.TunnelMessage) error {
	return w.replyFunc.Load().(ReplyFunc)(out)
}

func (w *Responser) ConsumeTunnelResponse() <-chan xnet.TunnelMessage {
	return w.msgs
}

func (w *Responser) ReplyFunc() ReplyFunc {
	return w.replyFunc.Load().(ReplyFunc)
}

func (w *Responser) UpdateReplyFunc(replyFunc ReplyFunc) {
	w.replyFunc.Store(replyFunc)
}
