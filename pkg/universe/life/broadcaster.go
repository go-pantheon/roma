package life

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	servicev1 "github.com/go-pantheon/roma/gen/api/server/gate/service/push/v1"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/go-pantheon/roma/pkg/universe/data"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type Broadcaster interface {
	Broadcast(wctx Context, mod climod.ModuleID, seq int32, obj int64, body proto.Message) error
	Multicast(wctx Context, mod climod.ModuleID, seq int32, obj int64, body proto.Message, uids ...int64) error
	Consume() <-chan *BroadcastMessage
	Send(msg *BroadcastMessage) error
}

var _ Broadcaster = (*WorkerBroadcaster)(nil)

// WorkerBroadcaster TODO
type WorkerBroadcaster struct {
	pusher *data.PushRepo
	msgs   chan *BroadcastMessage
}

func NewBroadcaster(pusher *data.PushRepo) Broadcaster {
	ep := &WorkerBroadcaster{
		pusher: pusher,
		msgs:   make(chan *BroadcastMessage, constants.WorkerReplySize),
	}
	return ep
}

func (w *WorkerBroadcaster) Broadcast(wctx Context, mod climod.ModuleID, seq int32, obj int64, body proto.Message) (err error) {
	bytes, err := proto.Marshal(body)
	if err != nil {
		err = errors.Wrapf(err, "[workerBroadcaster.Broadcast] proto.Marshal failed. mod<%d> seq<%d> obj<%d>", mod, seq, obj)
		return
	}

	w.msgs <- NewBroadcastMessage(true, nil, &servicev1.PushBody{
		Mod:  int32(mod),
		Seq:  seq,
		Obj:  obj,
		Data: bytes,
	})
	return
}

func (w *WorkerBroadcaster) Multicast(wctx Context, mod climod.ModuleID, seq int32, obj int64, body proto.Message, uids ...int64) (err error) {
	if len(uids) == 0 {
		err = errors.New("[workerBroadcaster.Multicast] uids is empty")
		return
	}

	bytes, err := proto.Marshal(body)
	if err != nil {
		err = errors.Wrapf(err, "[workerBroadcaster.Broadcast] proto.Marshal failed. mod<%d> seq<%d> obj<%d>", mod, seq, obj)
		return
	}

	w.msgs <- NewBroadcastMessage(false, uids, &servicev1.PushBody{
		Mod:  int32(mod),
		Seq:  seq,
		Obj:  obj,
		Data: bytes,
	})
	return
}

func (w *WorkerBroadcaster) Consume() <-chan *BroadcastMessage {
	return w.msgs
}

func (w *WorkerBroadcaster) Send(msg *BroadcastMessage) error {
	if msg.all {
		return w.pusher.Broadcast(context.Background(), msg.msg)
	}

	if len(msg.uids) > 16 {
		return w.pusher.Multicast(context.Background(), msg.uids, msg.msg)
	}

	for _, uid := range msg.uids {
		if err := w.pusher.Push(context.Background(), uid, msg.msg); err != nil {
			log.Errorf("push failed. uid<%d> err<%v>", uid, err)
		}
	}
	return nil
}

type BroadcastMessage struct {
	all  bool
	uids []int64
	msg  *servicev1.PushBody
}

func NewBroadcastMessage(all bool, uids []int64, msg *servicev1.PushBody) *BroadcastMessage {
	return &BroadcastMessage{
		all:  all,
		uids: uids,
		msg:  msg,
	}
}
