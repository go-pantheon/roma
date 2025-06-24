package core

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	"github.com/go-pantheon/roma/app/room/internal/client/self"
	"github.com/go-pantheon/roma/pkg/universe/data"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

type Manager struct {
	*life.Manager

	pusher *data.PushRepo
}

func NewManager(logger log.Logger, rt *self.SelfRouteTable, roomDo *domain.RoomDomain, pusher *data.PushRepo) (*Manager, func()) {
	newPersist := func(ctx context.Context, id int64, allowBorn bool) (hold life.Persistent, born bool, err error) {
		return newRoomPersister(ctx, roomDo, id, allowBorn)
	}

	lifeMgr, stopFunc := life.NewManager(logger, rt, newContext, newPersist)

	m := &Manager{
		Manager: lifeMgr,
		pusher:  pusher,
	}
	return m, func() {
		stopFunc()
	}
}

type eventFunc func(wctx Context, args ...int64) (err error)

func (m *Manager) SecondTickRegister(f func(ctx Context) error) {
	m.Manager.SecondTickRegister(func(ctx life.Context) error {
		return f(ctx.(Context))
	})
}

func (m *Manager) MinuteTickRegister(f func(ctx Context) error) {
	m.Manager.MinuteTickRegister(func(ctx life.Context) error {
		return f(ctx.(Context))
	})
}

func (m *Manager) CustomEventRegister(e life.WorkerEventType, f eventFunc) {
	m.Manager.CustomEventRegister(e, func(wctx life.Context, args ...int64) (err error) {
		return f(wctx.(Context), args...)
	})
}

func (m *Manager) OnLoadEventRegister(f func(ctx Context) error) {
	m.Manager.OnLoadEventRegister(func(ctx life.Context) error {
		return f(ctx.(Context))
	})
}

func (m *Manager) OnCreatedEventRegister(f func(ctx Context) error) {
	m.Manager.OnCreatedEventRegister(func(ctx life.Context) error {
		return f(ctx.(Context))
	})
}

func (m *Manager) Pusher() *data.PushRepo {
	return m.pusher
}

func (m *Manager) ExecuteAppEvent(ctx context.Context, oid int64, f life.EventFunc) error {
	w, err := m.Worker(ctx, oid, NewReplier(adminReplyFunc), life.NewBroadcaster(m.Pusher()))
	if err != nil {
		return err
	}

	return w.ExecuteEvent(newContext(ctx, w), f)
}
