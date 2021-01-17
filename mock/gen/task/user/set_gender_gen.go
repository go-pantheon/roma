// Code generated by gen-mercury. DO NOT EDIT.

package user

import (
	"github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	"github.com/vulcan-frame/vulcan-game/gen/api/client/sequence"
	"github.com/vulcan-frame/vulcan-game/mock/internal/task"
	"reflect"
)

var _ task.Taskable = (*SetGenderTask)(nil)

// SetGenderTask Set gender
type SetGenderTask struct {
	*task.CommonTask
}

func NewSetGenderTask(cs *climsg.CSSetGender, assert task.AssertFunc) *SetGenderTask {
	common := task.NewCommonTask(
		task.TypeCommon,
		climod.ModuleID_User,
		int32(cliseq.UserSeq_SetGender),
		reflect.TypeOf(climsg.SCSetGender{}),
		cs,
		assert,
	)
	o := &SetGenderTask{
		CommonTask: common,
	}
	return o
}
