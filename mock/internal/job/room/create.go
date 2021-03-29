package room

import (
	"github.com/pkg/errors"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/mock/gen/task/room"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job"
	"google.golang.org/protobuf/proto"
)

func NewCreateRoom() *job.Job {
	j := &job.Job{
		T: job.TypeRoom,
	}

	j.Tasks = append(j.Tasks, room.NewCreateRoomTask(&climsg.CSCreateRoom{
		RoomType: 1,
	}, assertCreateRoom))
	return j
}

func assertCreateRoom(ctx *base.Context, cs, sc proto.Message) (done bool, err error) {
	p, ok := sc.(*climsg.SCCreateRoom)
	if !ok {
		return
	}

	if p.Code != 1 {
		err = errors.Errorf("SCCreateRoom failed. code=%d", p.Code)
		return
	}
	done = true
	return
}
