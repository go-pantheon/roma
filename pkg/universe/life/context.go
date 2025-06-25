package life

import (
	"context"
	"sync"
	"time"

	"github.com/go-pantheon/fabrica-kit/xcontext"
)

type Context interface {
	context.Context
	EventManageable
	Responsive

	Now() time.Time
	ClientIP() string

	UID() int64
	OID() int64

	UnsafeObject() any
	Snapshot() VersionProto

	ChangedModules() (modules []ModuleKey, immediately bool)
	Changed(modules ...ModuleKey)
	ChangedImmediately(modules ...ModuleKey)
}

var _ Context = (*workerContext)(nil)

type workerContext struct {
	sync.Once
	context.Context
	EventManageable
	Responsive

	persister Persistent
	ctime     time.Time
	clientIP  string
	uid       int64

	changedModules     map[ModuleKey]struct{}
	changedImmediately bool
}

func NewContext(ctx context.Context, w *Worker) Context {
	c := &workerContext{
		Context:         ctx,
		EventManageable: w,
		Responsive:      w.Responsive,
		persister:       w.persistManager.Persister(),
		clientIP:        w.ClientIP(),
	}

	c.changedModules = make(map[ModuleKey]struct{}, len(c.persister.ModuleKeys()))

	if uid, err := xcontext.UID(ctx); err == nil {
		c.uid = uid
	}

	return c
}

func (c *workerContext) ChangedModules() (modules []ModuleKey, immediately bool) {
	defer func() {
		c.changedModules = make(map[ModuleKey]struct{}, len(c.persister.ModuleKeys()))
		c.changedImmediately = false
	}()

	modules = make([]ModuleKey, 0, len(c.changedModules))

	for mod := range c.changedModules {
		modules = append(modules, mod)
	}

	return modules, c.changedImmediately
}

func (c *workerContext) Changed(modules ...ModuleKey) {
	if len(modules) == 0 {
		for _, mod := range c.persister.ModuleKeys() {
			c.changedModules[mod] = struct{}{}
		}

		return
	}

	for _, mod := range modules {
		c.changedModules[mod] = struct{}{}
	}
}

func (c *workerContext) ChangedImmediately(modules ...ModuleKey) {
	c.changedImmediately = true

	c.Changed(modules...)
}

func (c *workerContext) Now() time.Time {
	c.Once.Do(func() {
		c.ctime = time.Now()
	})

	return c.ctime
}

func (c *workerContext) ClientIP() string {
	return c.clientIP
}

func (c *workerContext) UID() int64 {
	return c.uid
}

func (c *workerContext) OID() int64 {
	return c.persister.ID()
}

func (c *workerContext) Snapshot() VersionProto {
	return c.persister.Snapshot()
}

func (c *workerContext) UnsafeObject() any {
	return c.persister.UnsafeObject()
}
