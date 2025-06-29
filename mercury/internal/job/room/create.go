package room

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/room"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/job"
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

func assertCreateRoom(ctx core.Worker, cs, sc proto.Message) (err error) {
	return nil
}
