package storage

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/dev"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"google.golang.org/protobuf/proto"
)

func NewStorageJob() *job.Job {
	j := &job.Job{
		T: job.TypeStorage,
	}

	j.Tasks = append(j.Tasks, dev.NewDevListTask(&climsg.CSDevList{}, assert))

	j.Tasks = append(j.Tasks, dev.NewDevExecuteTask(&climsg.CSDevExecute{
		Mod:  "Storage",
		Cmd:  "Add Item",
		Args: map[string]string{"ItemID": "1", "Amount": "100"},
	}, assert))

	j.Tasks = append(j.Tasks, dev.NewDevExecuteTask(&climsg.CSDevExecute{
		Mod:  "Storage",
		Cmd:  "Clear",
		Args: map[string]string{},
	}, assert))

	return j
}

func assert(ctx core.Worker, cs, sc proto.Message) (err error) {
	return nil
}
