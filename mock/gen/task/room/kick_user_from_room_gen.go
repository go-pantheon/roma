// Code generated by gen-mock. DO NOT EDIT.

package room

import (
	"github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	"github.com/vulcan-frame/vulcan-game/gen/api/client/sequence"
	"github.com/vulcan-frame/vulcan-game/mock/internal/task"
	"reflect"
)

var _ task.Taskable = (*KickUserFromRoomTask)(nil)

// KickUserFromRoomTask Kick user from room
type KickUserFromRoomTask struct {
	*task.CommonTask
}

func NewKickUserFromRoomTask(cs *climsg.CSKickUserFromRoom, assert task.AssertFunc) *KickUserFromRoomTask {
	common := task.NewCommonTask(
		task.TypeCommon,
		climod.ModuleID_Room,
		int32(cliseq.RoomSeq_KickUserFromRoom),
		reflect.TypeOf(climsg.SCKickUserFromRoom{}),
		cs,
		assert,
	)
	o := &KickUserFromRoomTask{
		CommonTask: common,
	}
	return o
}
