package core

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate/domain"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/data"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/life"
)

type Manager struct {
	*life.Manager

	pusher *data.PushRepo
}

func NewManager(logger log.Logger, roomDo *domain.RoomDomain, pusher *data.PushRepo) (*Manager, func()) {
	newPersist := func(ctx context.Context, id int64, allowBorn bool) (hold life.Persistent, born bool, err error) {
		return newRoomPersister(ctx, roomDo, id, allowBorn)
	}

	lifeMgr, stopFunc := life.NewManager(logger, newContext, newPersist)

	m := &Manager{
		Manager: lifeMgr,
		pusher:  pusher,
	}
	return m, stopFunc
}

type eventFunc func(wctx Context, args ...int64) (err error)

func (m *Manager) SecondTickRegister(f func(ctx Context)) {
	m.Manager.SecondTickRegister(func(ctx life.Context) {
		f(ctx.(Context))
	})
}

func (m *Manager) MinuteTickRegister(f func(ctx Context)) {
	m.Manager.MinuteTickRegister(func(ctx life.Context) {
		f(ctx.(Context))
	})
}

func (m *Manager) CustomEventRegister(e life.WorkerEventType, f eventFunc) {
	m.Manager.CustomEventRegister(e, func(wctx life.Context, args ...int64) (err error) {
		return f(wctx.(Context), args...)
	})
}

func (m *Manager) OnLoadEventRegister(f func(ctx Context)) {
	m.Manager.OnLoadEventRegister(func(ctx life.Context) {
		f(ctx.(Context))
	})
}

func (m *Manager) OnCreatedEventRegister(f func(ctx Context)) {
	m.Manager.OnCreatedEventRegister(func(ctx life.Context) {
		f(ctx.(Context))
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
