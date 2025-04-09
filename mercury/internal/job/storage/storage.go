package storage

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/dev"
	"github.com/go-pantheon/roma/mercury/internal/base"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"google.golang.org/protobuf/proto"
)

func NewStorageJob() *job.Job {
	j := &job.Job{
		T: job.TypeStorage,
	}

	j.Tasks = append(j.Tasks, dev.NewDevListTask(&climsg.CSDevList{}, assert))
	j.Tasks = append(j.Tasks, dev.NewDevExecuteTask(&climsg.CSDevExecute{
		Mod:  "storage",
		Cmd:  "add",
		Args: map[string]string{"item_id": "1", "num": "100"},
	}, assert))
	j.Tasks = append(j.Tasks, dev.NewDevExecuteTask(&climsg.CSDevExecute{
		Mod:  "storage",
		Cmd:  "clear",
		Args: map[string]string{},
	}, assert))
	return j
}

func assert(ctx *base.Context, cs, sc proto.Message) (done bool, err error) {
	return true, nil
}
