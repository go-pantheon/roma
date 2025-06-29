package room

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/room"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"google.golang.org/protobuf/proto"
)

func NewGmList() *job.Job {
	j := &job.Job{
		T: job.TypeDev,
	}

	j.Tasks = append(j.Tasks, room.NewRoomListTask(&climsg.CSRoomList{}, assertGmList))

	return j
}

func assertGmList(ctx core.Worker, cs, sc proto.Message) (err error) {
	return nil
}
