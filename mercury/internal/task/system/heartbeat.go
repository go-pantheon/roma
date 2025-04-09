package system

import (
	"reflect"

	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	"github.com/go-pantheon/roma/mercury/internal/task"
)

var _ task.Taskable = (*HeartbeatTask)(nil)

type HeartbeatTask struct {
	*task.CommonTask
}

func NewHeartbeatTask(assert task.AssertFunc) *HeartbeatTask {
	common := task.NewCommonTask(
		task.TypeHeartbeat,
		climod.ModuleID_System,
		int32(cliseq.SystemSeq_Heartbeat),
		reflect.TypeOf(climsg.SCHeartbeat{}),
		nil,
		assert,
	)
	o := &HeartbeatTask{
		CommonTask: common,
	}
	return o
}
