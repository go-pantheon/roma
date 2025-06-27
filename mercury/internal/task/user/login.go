package user

import (
	"errors"

	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	clipkt "github.com/go-pantheon/roma/gen/api/client/packet"
	"github.com/go-pantheon/roma/mercury/gen/task/user"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/task"
)

var _ task.Taskable = (*LoginTask)(nil)

type LoginTask struct {
	*user.LoginTask
}

func NewLoginTask(cs *climsg.CSLogin, assert task.AssertFunc) *LoginTask {
	return &LoginTask{
		LoginTask: user.NewLoginTask(cs, assert),
	}
}

func (t *LoginTask) Receive(ctx *core.Context, sc *clipkt.Packet) (out *clipkt.Packet, done bool, err error) {
	if out, done, err = t.CommonTask.Receive(ctx, sc); err != nil || out != nil {
		return
	}

	b, err := t.UnmarshalSC(sc)
	if err != nil {
		return nil, false, err
	}

	body, ok := b.(*climsg.SCLogin)
	if !ok {
		return nil, false, errors.New("invalid sc message")
	}

	if done, err = t.Assert(ctx, t.CS, body); err != nil {
		return nil, false, err
	}

	ctx.Manager.SetClientUser(body.User)

	return nil, done, nil
}
