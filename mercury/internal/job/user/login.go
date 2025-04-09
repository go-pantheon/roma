package user

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/internal/base"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/go-pantheon/roma/mercury/internal/task/user"
	"github.com/pkg/errors"
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
