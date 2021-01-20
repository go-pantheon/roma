package storage

import (
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/mock/gen/task/dev"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job"
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
