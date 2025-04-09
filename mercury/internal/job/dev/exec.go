package dev

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/dev"
	"github.com/go-pantheon/roma/mercury/internal/base"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/pkg/errors"
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
