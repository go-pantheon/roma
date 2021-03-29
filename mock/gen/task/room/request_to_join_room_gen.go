// Code generated by gen-mock. DO NOT EDIT.

package room

import (
	"github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	"github.com/vulcan-frame/vulcan-game/gen/api/client/sequence"
	"github.com/vulcan-frame/vulcan-game/mock/internal/task"
	"reflect"
)

var _ task.Taskable = (*RequestToJoinRoomTask)(nil)

// RequestToJoinRoomTask Request to join room
type RequestToJoinRoomTask struct {
	*task.CommonTask
}

func NewRequestToJoinRoomTask(cs *climsg.CSRequestToJoinRoom, assert task.AssertFunc) *RequestToJoinRoomTask {
	common := task.NewCommonTask(
		task.TypeCommon,
		climod.ModuleID_Room,
		int32(cliseq.RoomSeq_RequestToJoinRoom),
		reflect.TypeOf(climsg.SCRequestToJoinRoom{}),
		cs,
		assert,
	)
	o := &RequestToJoinRoomTask{
		CommonTask: common,
	}
	return o
}
