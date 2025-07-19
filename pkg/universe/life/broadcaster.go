package life

import (
	"context"
	"sync"

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
	ConsumeBroadcast() <-chan *BroadcastMessage
	ExecuteBroadcast(ctx context.Context, msg *BroadcastMessage) error
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
		return errors.Wrapf(err, "[workerBroadcaster.Broadcast] proto.Marshal failed. %d-%d %+v", mod, seq, obj)
	}

	w.msgs <- newBroadcastMessage(true, nil, buildPushBody(mod, seq, obj, bytes))

	return nil
}

func (w *WorkerBroadcaster) Multicast(wctx Context, mod climod.ModuleID, seq int32, obj int64, body proto.Message, uids ...int64) (err error) {
	if len(uids) == 0 {
		return errors.New("[workerBroadcaster.Multicast] uids is empty")
	}

	bytes, err := proto.Marshal(body)
	if err != nil {
		return errors.Wrapf(err, "[workerBroadcaster.Broadcast] proto.Marshal failed. %d-%d %+v", mod, seq, obj)
	}

	w.msgs <- newBroadcastMessage(false, uids, buildPushBody(mod, seq, obj, bytes))

	return nil
}

func (w *WorkerBroadcaster) ConsumeBroadcast() <-chan *BroadcastMessage {
	return w.msgs
}

func (w *WorkerBroadcaster) ExecuteBroadcast(ctx context.Context, msg *BroadcastMessage) error {
	defer putBroadcastMessage(msg)

	if msg.all {
		return w.pusher.Broadcast(ctx, msg.out)
	}

	return w.pusher.Multicast(ctx, msg.uids, msg.out)
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
