package user

import (
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	clipkt "github.com/vulcan-frame/vulcan-game/gen/api/client/packet"
	"github.com/vulcan-frame/vulcan-game/mock/gen/task/user"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/task"
)

var _ task.Taskable = (*LoginTask)(nil)

type LoginTask struct {
	*user.LoginTask
}

func NewLoginTask(cs *climsg.CSLogin, assert task.AssertFunc) *LoginTask {
	o := &LoginTask{
		LoginTask: user.NewLoginTask(cs, assert),
	}
	return o
}

func (t *LoginTask) Receive(ctx *base.Context, sc *clipkt.Packet) (out *clipkt.Packet, done bool, err error) {
	if out, done, err = t.CommonTask.Receive(ctx, sc); err != nil || out != nil {
		return
	}

	b, err := t.UnmarshalSC(sc)
	if err != nil {
		return
	}

	body, _ := b.(*climsg.SCLogin)
	if done, err = t.Assert(ctx, t.CS, body); err != nil {
		return
	}

	ctx.Manager.SetClientUser(body.User)
	return
}
