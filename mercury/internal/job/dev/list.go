package dev

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/dev"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"google.golang.org/protobuf/proto"
)

func NewDevListJob() *job.Job {
	j := &job.Job{
		T: job.TypeDev,
	}

	j.Tasks = append(j.Tasks, dev.NewDevListTask(&climsg.CSDevList{}, assertDevList))

	return j
}

func assertDevList(ctx core.Worker, cs, sc proto.Message) (err error) {
	return nil
}
