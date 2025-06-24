package core

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/client/self"
	"github.com/go-pantheon/roma/pkg/universe/data"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

type Manager struct {
	*life.Manager

	pusher *data.PushRepo
}

func NewManager(logger log.Logger, rt *self.SelfRouteTable, userDo *domain.UserDomain, pusher *data.PushRepo) (*Manager, func()) {
	newPersister := func(ctx context.Context, uid int64, allowBorn bool) (persister life.Persistent, born bool, err error) {
		return newUserPersister(ctx, userDo, uid, allowBorn)
	}

	lifeMgr, stopFunc := life.NewManager(logger, rt, newContext, newPersister)

	m := &Manager{
		Manager: lifeMgr,
		pusher:  pusher,
	}

	return m, func() {
		stopFunc()
	}
}

type eventFunc func(wctx Context, arg *life.EventArg) (err error)

func (m *Manager) RegisterEvent(et life.WorkerEventType, f eventFunc) {
	m.Manager.RegisterCustomEvent(et, func(wctx life.Context, arg *life.EventArg) (err error) {
		return f(wctx.(Context), arg)
	})
}

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

func (m *Manager) Pusher() *data.PushRepo {
	return m.pusher
}

func (m *Manager) ExecuteEvent(ctx context.Context, uid int64, f life.EventFunc) (err error) {
	w, err := m.Worker(ctx, uid, NewReplier(EmptyReplyFunc), life.NewBroadcaster(m.Pusher()))
	if err != nil {
		return err
	}

	defer func() {
		// check the status before stop the worker to avoid incorrect stop the worker that will be connected to gate when the admin execute the event
		if life.IsGateStatus(w.Status()) {
			return
		}

		if stoperr := w.Stop(ctx); stoperr != nil {
			err = errors.Join(err, stoperr)
		}
	}()

	return w.ExecuteEvent(newContext(ctx, w), f)
}
