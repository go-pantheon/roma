package system

import (
	"time"

	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/system"
	"github.com/go-pantheon/roma/mercury/internal/base"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func NewHeartbeatJob() *job.Job {
	j := &job.Job{
		T: job.TypeSystem,
	}

	j.Tasks = append(j.Tasks, NewHeartbeatTask())

	return j
}

func assertHeartBeatExec(ctx *base.Context, cs, sc proto.Message) (done bool, err error) {
	p, ok := sc.(*climsg.SCHeartbeat)
	if !ok {
		return
	}

	if p.Code != 1 {
		err = errors.Errorf("SCHeartbeat failed. code=%d", p.Code)
		return
	}
	done = true
	return
}

func NewHeartbeatTask() *system.HeartbeatTask {
	return system.NewHeartbeatTask(&climsg.CSHeartbeat{
		ClientTime: time.Now().Unix(),
	}, assertHeartBeatExec)
}
