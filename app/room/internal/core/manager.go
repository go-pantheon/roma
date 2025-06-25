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
}

func NewManager(logger log.Logger, rt *self.SelfRouteTable, roomDo *domain.RoomDomain, pusher *data.PushRepo) (*Manager, func()) {
	newPersist := func(ctx context.Context, id int64, allowBorn bool) (hold life.Persistent, born bool, err error) {
		return newRoomPersister(ctx, roomDo, id, allowBorn)
	}

	lifeMgr, stopFunc := life.NewManager(logger, rt, pusher, newContext, newPersist)

	m := &Manager{
		Manager: lifeMgr,
	}

	return m, func() {
		stopFunc()
	}
}

type eventFunc func(wctx Context, arg *life.EventArg) (err error)

func (m *Manager) RegisterSecondTick(f func(ctx Context) error) {
	m.Manager.RegisterSecondTick(func(ctx life.Context) error {
		return f(ctx.(Context))
	})
}

func (m *Manager) RegisterMinuteTick(f func(ctx Context) error) {
	m.Manager.RegisterMinuteTick(func(ctx life.Context) error {
		return f(ctx.(Context))
	})
}

func (m *Manager) RegisterEvent(e life.WorkerEventType, f eventFunc) {
	m.Manager.RegisterCustomEvent(e, func(wctx life.Context, arg *life.EventArg) (err error) {
		return f(wctx.(Context), arg)
	})
}

func (m *Manager) RegisterOnLoadEvent(f func(ctx Context) error) {
	m.Manager.RegisterOnLoadEvent(func(ctx life.Context) error {
		return f(ctx.(Context))
	})
}

func (m *Manager) RegisterOnCreatedEvent(f func(ctx Context) error) {
	m.Manager.RegisterOnCreatedEvent(func(ctx life.Context) error {
		return f(ctx.(Context))
	})
}

func (m *Manager) ExecuteEvent(ctx context.Context, oid int64, f life.EventFunc) error {
	w, err := m.Worker(ctx, oid, NewResponser(mockResponseFunc), life.NewBroadcaster(m.Pusher()))
	if err != nil {
		return err
	}

	return w.ExecuteEvent(newContext(ctx, w), f)
}
