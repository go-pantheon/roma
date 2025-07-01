package core

import (
	"context"

	userobj "github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/object"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

type Context interface {
	life.Context

	User() *userobj.User
}

var _ Context = (*gameContext)(nil)

type gameContext struct {
	life.Context
}

func newContext(ctx context.Context, w *life.Worker) life.Context {
	c := &gameContext{
		Context: life.NewContext(ctx, w),
	}

	return c
}

func (w *gameContext) User() *userobj.User {
	return w.UnsafeObject().(*userobj.User)
}

func (w *gameContext) UID() int64 {
	return w.OID()
}
