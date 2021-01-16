package life

import (
	"context"
	"sync"
	"time"

	"github.com/vulcan-frame/vulcan-kit/xcontext"
	"google.golang.org/protobuf/proto"
)

type Context interface {
	context.Context
	EventManageable
	Replier

	UID() int64
	SetUID(uid int64)
	Now() time.Time
	ClientIP() string

	OID() int64
	UnsafeObject() interface{}
	ShowProto() proto.Message

	IsChanged() (changed bool, immediately bool)
	Changed()
	ChangedImmediately()
}

var _ Context = (*workerContext)(nil)

type workerContext struct {
	sync.Once
	context.Context
	EventManageable
	Replier

	persister Persistent
	ctime     time.Time
	clientIP  string
	uid       int64

	changed            bool
	changedImmediately bool
}

func NewContext(ctx context.Context, w *Worker) Context {
	c := &workerContext{
		Context:         ctx,
		EventManageable: w,
		Replier:         w.Replier,
		persister:       w.persistManager.Persister(),
		clientIP:        w.ClientIP(),
	}
	if uid, err := xcontext.UID(ctx); err == nil {
		c.SetUID(uid)
	}
	return c
}

func (w *workerContext) Now() time.Time {
	w.Once.Do(func() {
		w.ctime = time.Now()
	})
	return w.ctime
}

func (w *workerContext) ClientIP() string {
	return w.clientIP
}

func (w *workerContext) SetUID(uid int64) {
	w.uid = uid
}

func (w *workerContext) UID() int64 {
	return w.uid
}

func (w *workerContext) OID() int64 {
	return w.persister.ID()
}

func (w *workerContext) ShowProto() proto.Message {
	return w.persister.ShowProto()
}

func (w *workerContext) UnsafeObject() interface{} {
	return w.persister.UnsafeObject()
}

func (w *workerContext) IsChanged() (changed bool, immediately bool) {
	return w.changed, w.changedImmediately
}

func (w *workerContext) Changed() {
	w.changed = true
}

func (w *workerContext) ChangedImmediately() {
	w.changed = true
	w.changedImmediately = true
}
