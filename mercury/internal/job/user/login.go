package user

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/go-pantheon/roma/mercury/internal/task/user"
	"google.golang.org/protobuf/proto"
)

func NewLoginJob() *job.Job {
	j := &job.Job{
		T: job.TypeUserLogin,
	}

	j.Tasks = append(j.Tasks, user.NewLoginTask(&climsg.CSLogin{}, assertLogin))

	return j
}

func assertLogin(ctx core.Worker, cs, sc proto.Message) (err error) {
	return nil
}
