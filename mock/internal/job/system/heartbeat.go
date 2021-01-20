package system

import (
	"time"

	"github.com/pkg/errors"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/mock/gen/task/system"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job"
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
