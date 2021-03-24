package dev

import (
	"github.com/pkg/errors"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/mock/gen/task/dev"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job"
	"google.golang.org/protobuf/proto"
)

func NewDevExec() *job.Job {
	j := &job.Job{
		T: job.TypeDev,
	}

	j.Tasks = append(j.Tasks, dev.NewDevExecuteTask(&climsg.CSDevExecute{
		Mod:  "hero",
		Cmd:  "level_up",
		Args: map[string]string{"hero_id": "1", "target_level": "10"},
	}, assertDevExec))

	return j
}

func assertDevExec(ctx *base.Context, cs, sc proto.Message) (done bool, err error) {
	p, ok := sc.(*climsg.SCDevExecute)
	if !ok {
		return
	}

	if p.Code != 1 {
		err = errors.Errorf("SCDevExecute failed. code=%d", p.Code)
		return
	}
	done = true
	return
}
