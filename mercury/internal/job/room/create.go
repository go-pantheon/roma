package room

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/room"
	"github.com/go-pantheon/roma/mercury/internal/base"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/pkg/errors"
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
