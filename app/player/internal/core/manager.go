package core

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/domain"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/life"
)

type Manager struct {
	*life.Manager
}

func NewManager(logger log.Logger, userDo *domain.UserDomain) (*Manager, func()) {
	newPersister := func(ctx context.Context, uid int64, allowBorn bool) (persister life.Persistent, born bool, err error) {
		return newUserPersister(ctx, userDo, uid, allowBorn)
	}

	lifeMgr, stopFunc := life.NewManager(logger, newContext, newPersister)
	m := &Manager{
		Manager: lifeMgr,
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
