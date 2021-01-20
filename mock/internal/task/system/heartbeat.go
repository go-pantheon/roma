package system

import (
	"reflect"

	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	climod "github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	cliseq "github.com/vulcan-frame/vulcan-game/gen/api/client/sequence"
	"github.com/vulcan-frame/vulcan-game/mock/internal/task"
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
