package core

import (
	"sync/atomic"

	"github.com/go-kratos/kratos/log"
	jsoniter "github.com/json-iterator/go"
	climod "github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	clipkt "github.com/vulcan-frame/vulcan-game/gen/api/client/packet"
	intrav1 "github.com/vulcan-frame/vulcan-game/gen/api/server/room/intra/v1"
	"github.com/vulcan-frame/vulcan-game/gen/app/codec"
	"github.com/vulcan-frame/vulcan-game/gen/app/room/handler"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/constants"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

var _ life.Replier = (*Replier)(nil)

var adminReplyFunc = func(out proto.Message) error {
	reply, ok := out.(*intrav1.TunnelResponse)
	if !ok {
		log.Errorf("intrav1.TunnelResponse type conversion failed. out=%T", out)
		return nil
	}

	p := &clipkt.Packet{
		Mod:  reply.GetMod(),
		Seq:  reply.GetSeq(),
		Obj:  reply.GetObj(),
		Data: reply.GetData(),
	}

	body, _ := codec.UnmarshalSC(p)
	str, _ := jsoniter.MarshalToString(body)

	var tag string
	if codec.IsPushSC(climod.ModuleID(p.Mod), p.Seq) {
		tag = "Push"
	} else {
		tag = "Reply"
	}
	log.Debugf("[Admin] [%s] %d-%d oid=%d body=%s", tag, p.Mod, p.Seq, p.Obj, str)

	return nil
}

type Replier struct {
	msgs      chan proto.Message
	replyFunc atomic.Value
}

func NewReplier(replyFunc life.ReplyFunc) life.Replier {
	o := &Replier{
		msgs: make(chan proto.Message, constants.WorkerReplySize),
	}
	o.replyFunc.Store(replyFunc)
	return o
}

func (w *Replier) Reply(mod climod.ModuleID, seq int32, obj int64, sc proto.Message) error {
	msg, err := handler.NewRoomResponseProto(int32(mod), seq, obj, sc)
	if err != nil {
		return err
	}
	w.msgs <- msg
	return nil
}

func (w *Replier) ReplyImmediately(mod climod.ModuleID, seq int32, obj int64, sc proto.Message) error {
	msg, err := handler.NewRoomResponseProto(int32(mod), seq, obj, sc)
	if err != nil {
		return err
	}
	return w.ExecuteReply(msg)
}

func (w *Replier) ReplyBytes(mod climod.ModuleID, seq int32, obj int64, out []byte) error {
	w.msgs <- handler.NewRoomResponseProtoByData(int32(mod), seq, obj, out)
	return nil
}

func (w *Replier) ExecuteReply(out proto.Message) error {
	return w.replyFunc.Load().(life.ReplyFunc)(out)
}

func (w *Replier) ConsumeReplyMessage() <-chan proto.Message {
	return w.msgs
}

func (w *Replier) UpdateReplyFunc(replyFunc life.ReplyFunc) {
	w.replyFunc.Store(replyFunc)
}

func (w *Replier) GetReplyFunc() life.ReplyFunc {
	return w.replyFunc.Load().(life.ReplyFunc)
}
