package life

import (
	"context"
	"sync"
	"time"

	"github.com/go-pantheon/fabrica-kit/xcontext"
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

	ChangedModules() (modules []ModuleKey, immediately bool)
	Changed(modules ...ModuleKey)
	ChangedImmediately(modules ...ModuleKey)
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

	changedModules     []ModuleKey
	changedImmediately bool
}

const defaultChangedModulesSize = 16

func NewContext(ctx context.Context, w *Worker) Context {
	c := &workerContext{
		Context:         ctx,
		EventManageable: w,
		Replier:         w.Replier,
		persister:       w.persistManager.Persister(),
		clientIP:        w.ClientIP(),
		changedModules:  make([]ModuleKey, 0, defaultChangedModulesSize),
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

var emptyModules = []ModuleKey{}

func (w *workerContext) ChangedModules() (modules []ModuleKey, immediately bool) {
	if len(w.changedModules) == 0 {
		return emptyModules, false
	}

	modules = w.changedModules
	w.changedModules = make([]ModuleKey, 0, len(modules))

	immediately = w.changedImmediately
	w.changedImmediately = false

	return
}

func (w *workerContext) Changed(modules ...ModuleKey) {
	w.changedModules = append(w.changedModules, modules...)
}

func (w *workerContext) ChangedImmediately(modules ...ModuleKey) {
	w.changedImmediately = true
	w.changedModules = append(w.changedModules, modules...)
}
