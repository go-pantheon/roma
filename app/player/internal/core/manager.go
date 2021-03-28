package core

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/domain"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/data"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/life"
)

type Manager struct {
	*life.Manager

	pusher *data.PushRepo
}

func NewManager(logger log.Logger, userDo *domain.UserDomain, pusher *data.PushRepo) (*Manager, func()) {
	newPersister := func(ctx context.Context, uid int64, allowBorn bool) (persister life.Persistent, born bool, err error) {
		return newUserPersister(ctx, userDo, uid, allowBorn)
	}

	lifeMgr, stopFunc := life.NewManager(logger, newContext, newPersister)
	m := &Manager{
		Manager: lifeMgr,
		pusher:  pusher,
	}

	return m, stopFunc
}

type eventFunc func(wctx Context, args ...int64) (err error)

func (m *Manager) EventRegister(et life.WorkerEventType, f eventFunc) {
	m.Manager.CustomEventRegister(et, func(wctx life.Context, args ...int64) (err error) {
		return f(wctx.(Context), args...)
	})
}

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

func (m *Manager) ExecuteAppEvent(ctx context.Context, uid int64, f life.EventFunc) error {
	w, err := m.Worker(ctx, uid, NewReplier(EmptyReplyFunc), life.NewBroadcaster(m.Pusher()))
	if err != nil {
		return err
	}

	defer func() {
		// check the status before stop the worker, because the worker will be connected to gate when the admin execute the event
		if life.IsGateContext(ctx) {
			return
		}
		w.TriggerStop()
	}()

	return w.ExecuteEvent(newContext(ctx, w), f)
}
