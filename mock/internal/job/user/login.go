package user

import (
	"github.com/pkg/errors"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job"
	"github.com/vulcan-frame/vulcan-game/mock/internal/task/user"
	"google.golang.org/protobuf/proto"
)

func NewLoginJob() *job.Job {
	j := &job.Job{
		T: job.TypeUserLogin,
	}

	j.Tasks = append(j.Tasks, user.NewLoginTask(&climsg.CSLogin{}, assertLogin))

	return j
}

func assertLogin(ctx *base.Context, cs, sc proto.Message) (done bool, err error) {
	p, ok := sc.(*climsg.SCLogin)
	if !ok {
		return
	}

	if p.Code != 1 {
		err = errors.Errorf("SCLogin failed. code=%d", p.Code)
		return
	}

	done = true
	return
}
