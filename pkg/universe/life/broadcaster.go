package life

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	servicev1 "github.com/go-pantheon/roma/gen/api/server/broadcaster/service/push/v1"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/go-pantheon/roma/pkg/universe/data"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type Broadcastable interface {
	Broadcast(wctx Context, mod climod.ModuleID, seq int32, obj int64, body proto.Message) error
	Multicast(wctx Context, mod climod.ModuleID, seq int32, obj int64, body proto.Message, uids ...int64) error
	Consume() <-chan *BroadcastMessage
	ExecuteBroadcast(msg *BroadcastMessage) error
}

var _ Broadcastable = (*WorkerBroadcaster)(nil)

// WorkerBroadcaster TODO
type WorkerBroadcaster struct {
	pusher *data.PushRepo
	msgs   chan *BroadcastMessage
}

func NewBroadcaster(pusher *data.PushRepo) Broadcastable {
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

	w.msgs <- newBroadcastMessage(true, nil, buildPushBody(mod, seq, obj, bytes))

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

	w.msgs <- newBroadcastMessage(false, uids, buildPushBody(mod, seq, obj, bytes))

	return
}

func (w *WorkerBroadcaster) Consume() <-chan *BroadcastMessage {
	return w.msgs
}

func (w *WorkerBroadcaster) ExecuteBroadcast(msg *BroadcastMessage) error {
	defer putBroadcastMessage(msg)

	if msg.all {
		return w.pusher.Broadcast(context.Background(), msg.out)
	}

	if len(msg.uids) > 16 {
		return w.pusher.Multicast(context.Background(), msg.uids, msg.out)
	}

	for _, uid := range msg.uids {
		if err := w.pusher.Push(context.Background(), uid, msg.out); err != nil {
			log.Errorf("push failed. uid<%d> err<%v>", uid, err)
		}
	}

	return nil
}

type BroadcastMessage struct {
	all  bool
	uids []int64
	out  *servicev1.PushBody
}

func newBroadcastMessage(all bool, uids []int64, out *servicev1.PushBody) *BroadcastMessage {
	msg := getBroadcastMessage()

	msg.all = all
	msg.uids = uids
	msg.out = out

	return msg
}

func (b *BroadcastMessage) reset() {
	b.all = false
	b.uids = nil

	if b.out != nil {
		putPushBody(b.out)
		b.out = nil
	}
}

var (
	broadcastMessagePool = sync.Pool{
		New: func() any {
			return &BroadcastMessage{}
		},
	}

	pushBodyPool = sync.Pool{
		New: func() any {
			return &servicev1.PushBody{}
		},
	}
)

func getBroadcastMessage() *BroadcastMessage {
	return broadcastMessagePool.Get().(*BroadcastMessage)
}

func putBroadcastMessage(msg *BroadcastMessage) {
	msg.reset()
	broadcastMessagePool.Put(msg)
}

func buildPushBody(mod climod.ModuleID, seq int32, obj int64, data []byte) *servicev1.PushBody {
	body := pushBodyPool.Get().(*servicev1.PushBody)

	body.Mod = int32(mod)
	body.Seq = seq
	body.Obj = obj
	body.Data = data

	return body
}

func putPushBody(body *servicev1.PushBody) {
	body.Reset()
	pushBodyPool.Put(body)
}
