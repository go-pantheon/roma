package life

import (
	"sync/atomic"

	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"google.golang.org/protobuf/proto"
)

type SendFunc = func(msg xnet.TunnelMessage) error

type Responsive interface {
	Push(mod int32, seq int32, obj int64, sc proto.Message) error
	PushImmediately(mod int32, seq int32, obj int64, sc proto.Message) error
	Reply(msg xnet.TunnelMessage)
	ConsumeTunnelResponse() <-chan xnet.TunnelMessage
	ExecuteSend(msg xnet.TunnelMessage) error
	SendFunc() SendFunc
	UpdateSendFunc(sendFunc SendFunc)
	SCIndex() int32
}

type BuildTunnelResponseFunc = func(mod int32, seq int32, obj int64, msg proto.Message) (xnet.TunnelMessage, error)

var _ Responsive = (*Responser)(nil)

type Responser struct {
	msgs     chan xnet.TunnelMessage
	sendFunc atomic.Value

	scIndexFunc          func() int32
	buildPushMessageFunc BuildTunnelResponseFunc
}

func NewResponser(sendFunc SendFunc, buildPushMessageFunc BuildTunnelResponseFunc) *Responser {
	o := &Responser{
		msgs:                 make(chan xnet.TunnelMessage, constants.WorkerReplySize),
		sendFunc:             atomic.Value{},
		buildPushMessageFunc: buildPushMessageFunc,
	}

	o.sendFunc.Store(sendFunc)

	return o
}

func (w *Responser) Reply(msg xnet.TunnelMessage) {
	w.msgs <- msg
}

func (w *Responser) Push(mod int32, seq int32, obj int64, sc proto.Message) error {
	msg, err := w.buildPushMessageFunc(mod, seq, obj, sc)
	if err != nil {
		return err
	}

	w.msgs <- msg

	return nil
}

func (w *Responser) PushImmediately(mod int32, seq int32, obj int64, sc proto.Message) error {
	msg, err := w.buildPushMessageFunc(mod, seq, obj, sc)
	if err != nil {
		return err
	}

	return w.sendFunc.Load().(SendFunc)(msg)
}

func (w *Responser) ExecuteSend(out xnet.TunnelMessage) error {
	return w.sendFunc.Load().(SendFunc)(out)
}

func (w *Responser) ConsumeTunnelResponse() <-chan xnet.TunnelMessage {
	return w.msgs
}

func (w *Responser) SendFunc() SendFunc {
	return w.sendFunc.Load().(SendFunc)
}

func (w *Responser) UpdateSendFunc(sendFunc SendFunc) {
	w.sendFunc.Store(sendFunc)
}

func (w *Responser) SCIndex() int32 {
	return w.scIndexFunc()
}

func (w *Responser) SetSCIndexFunc(scIndexFunc func() int32) {
	w.scIndexFunc = scIndexFunc
}
