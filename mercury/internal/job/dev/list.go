package dev

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/dev"
	"github.com/go-pantheon/roma/mercury/internal/base"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/pkg/errors"
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
