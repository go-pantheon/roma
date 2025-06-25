package life

import (
	"sync/atomic"

	"github.com/go-pantheon/roma/pkg/universe/constants"
	"google.golang.org/protobuf/proto"
)

type ReplyFunc = func(msg proto.Message) error

type Responsive interface {
	Reply(mod int32, seq int32, obj int64, sc proto.Message) error
	ReplyImmediately(mod int32, seq int32, obj int64, sc proto.Message) error
	ReplyBytes(mod int32, seq int32, obj int64, body []byte)
	ConsumeReplyMessage() <-chan proto.Message
	ExecuteReply(msg proto.Message) error
	ReplyFunc() ReplyFunc
	UpdateReplyFunc(replyFunc ReplyFunc)
}

type BuildMsgByProtoFunc = func(mod int32, seq int32, obj int64, sc proto.Message) (proto.Message, error)
type BuildMsgByBytesFunc = func(mod int32, seq int32, obj int64, body []byte) proto.Message

var _ Responsive = (*Responser)(nil)

type Responser struct {
	msgs      chan proto.Message
	replyFunc atomic.Value

	buildMsgByProtoFunc BuildMsgByProtoFunc
	buildMsgByBytesFunc BuildMsgByBytesFunc
}

func NewResponser(replyFunc ReplyFunc, buildMsgByProtoFunc BuildMsgByProtoFunc, buildMsgByBytesFunc BuildMsgByBytesFunc) *Responser {
	o := &Responser{
		msgs:                make(chan proto.Message, constants.WorkerReplySize),
		replyFunc:           atomic.Value{},
		buildMsgByProtoFunc: buildMsgByProtoFunc,
		buildMsgByBytesFunc: buildMsgByBytesFunc,
	}

	o.replyFunc.Store(replyFunc)

	return o
}

func (w *Responser) Reply(mod int32, seq int32, obj int64, sc proto.Message) error {
	msg, err := w.buildMsgByProtoFunc(mod, seq, obj, sc)
	if err != nil {
		return err
	}

	w.msgs <- msg

	return nil
}

func (w *Responser) ReplyImmediately(mod int32, seq int32, obj int64, sc proto.Message) error {
	msg, err := w.buildMsgByProtoFunc(mod, seq, obj, sc)
	if err != nil {
		return err
	}

	return w.replyFunc.Load().(ReplyFunc)(msg)
}

func (w *Responser) ReplyBytes(mod int32, seq int32, obj int64, body []byte) {
	w.msgs <- w.buildMsgByBytesFunc(mod, seq, obj, body)
}

func (w *Responser) ExecuteReply(out proto.Message) error {
	return w.replyFunc.Load().(ReplyFunc)(out)
}

func (w *Responser) ConsumeReplyMessage() <-chan proto.Message {
	return w.msgs
}

func (w *Responser) ReplyFunc() ReplyFunc {
	return w.replyFunc.Load().(ReplyFunc)
}

func (w *Responser) UpdateReplyFunc(replyFunc ReplyFunc) {
	w.replyFunc.Store(replyFunc)
}
