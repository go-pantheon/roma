package core

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	"github.com/go-pantheon/roma/pkg/universe/data"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

type Manager struct {
	*life.Manager

	pusher *data.PushRepo
}

func NewManager(logger log.Logger, userDo *domain.UserDomain, pusher *data.PushRepo) (*Manager, func()) {
	newPersister := func(ctx context.Context, uid int64, sid int64, allowBorn bool) (persister life.Persistent, born bool, err error) {
		return newUserPersister(ctx, userDo, uid, sid, allowBorn)
	}

	lifeMgr, stopFunc := life.NewManager(logger, newContext, newPersister)
	m := &Manager{
		Manager: lifeMgr,
		pusher:  pusher,
	}

	return m, func() {
		stopFunc()
	}
}

type eventFunc func(wctx Context, args ...int64) (err error)

func (m *Manager) EventRegister(et life.WorkerEventType, f eventFunc) {
	m.Manager.CustomEventRegister(et, func(wctx life.Context, args ...int64) (err error) {
		return f(wctx.(Context), args...)
	})
}

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

func (m *Manager) ExecuteAppEvent(ctx context.Context, uid int64, sid int64, f life.EventFunc) (err error) {
	w, err := m.Worker(ctx, uid, sid, NewReplier(EmptyReplyFunc), life.NewBroadcaster(m.Pusher()))
	if err != nil {
		return err
	}

	defer func() {
		// check the status before stop the worker, because the worker will be connected to gate when the admin execute the event
		if life.IsGateContext(ctx) {
			return
		}

		if stoperr := w.Stop(ctx); stoperr != nil {
			err = errors.Join(err, stoperr)
		}
	}()

	return w.ExecuteEvent(newContext(ctx, w), f)
}
