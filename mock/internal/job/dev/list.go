package dev

import (
	"github.com/pkg/errors"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/mock/gen/task/dev"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job"
	"google.golang.org/protobuf/proto"
)

func NewDevListJob() *job.Job {
	j := &job.Job{
		T: job.TypeDev,
	}

	j.Tasks = append(j.Tasks, dev.NewDevListTask(&climsg.CSDevList{}, assertDevList))

	return j
}

func assertDevList(ctx *base.Context, cs, sc proto.Message) (done bool, err error) {
	p, ok := sc.(*climsg.SCDevList)
	if !ok {
		return
	}

	if p.Code != 1 {
		err = errors.Errorf("SCDevList failed. code=%d", p.Code)
		return
	}
	done = true
	return
}
