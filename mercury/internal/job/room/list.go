package room

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/room"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func NewGmList() *job.Job {
	j := &job.Job{
		T: job.TypeDev,
	}

	j.Tasks = append(j.Tasks, room.NewRoomListTask(&climsg.CSRoomList{}, assertGmList))

	return j
}

func assertGmList(ctx *core.Context, cs, sc proto.Message) (done bool, err error) {
	p, ok := sc.(*climsg.SCRoomList)
	if !ok {
		return false, errors.New("invalid sc message")
	}

	if p.Code != 1 {
		return false, errors.Errorf("SCRoomList failed. code=%d", p.Code)
	}

	return true, nil
}
