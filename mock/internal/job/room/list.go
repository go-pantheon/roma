package room

import (
	"github.com/pkg/errors"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/mock/gen/task/room"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job"
	"google.golang.org/protobuf/proto"
)

func NewGmList() *job.Job {
	j := &job.Job{
		T: job.TypeDev,
	}

	j.Tasks = append(j.Tasks, room.NewRoomListTask(&climsg.CSRoomList{}, assertGmList))

	return j
}

func assertGmList(ctx *base.Context, cs, sc proto.Message) (done bool, err error) {
	p, ok := sc.(*climsg.SCRoomList)
	if !ok {
		return
	}

	if p.Code != 1 {
		err = errors.Errorf("SCRoomList failed. code=%d", p.Code)
		return
	}
	done = true
	return
}
