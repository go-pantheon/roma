package user

import (
	"github.com/go-pantheon/fabrica-util/errors"
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

func (t *LoginTask) Receive(mgr core.Worker, sc *clipkt.Packet) (err error) {
	if err = t.CommonTask.Receive(mgr, sc); err != nil {
		return
	}

	b, err := t.UnmarshalSC(sc)
	if err != nil {
		return err
	}

	body, ok := b.(*climsg.SCLogin)
	if !ok {
		return errors.New("not SCLogin")
	}

	if err = t.Assert(mgr, t.CS, body); err != nil {
		return err
	}

	mgr.SetClientUser(body.User)

	return nil
}
