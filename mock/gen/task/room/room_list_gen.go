// Code generated by gen-mock. DO NOT EDIT.

package room

import (
	"github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	"github.com/vulcan-frame/vulcan-game/gen/api/client/sequence"
	"github.com/vulcan-frame/vulcan-game/mock/internal/task"
	"reflect"
)

var _ task.Taskable = (*RoomListTask)(nil)

// RoomListTask Room list
type RoomListTask struct {
	*task.CommonTask
}

func NewRoomListTask(cs *climsg.CSRoomList, assert task.AssertFunc) *RoomListTask {
	common := task.NewCommonTask(
		task.TypeCommon,
		climod.ModuleID_Room,
		int32(cliseq.RoomSeq_RoomList),
		reflect.TypeOf(climsg.SCRoomList{}),
		cs,
		assert,
	)
	o := &RoomListTask{
		CommonTask: common,
	}
	return o
}
