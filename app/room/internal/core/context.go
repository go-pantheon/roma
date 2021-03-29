package core

import (
	"context"

	roomobj "github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate/domain/object"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/life"
)

type Context interface {
	life.Context

	Room() *roomobj.Room
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

func (w *gameContext) Room() *roomobj.Room {
	return w.UnsafeObject().(*roomobj.Room)
}
