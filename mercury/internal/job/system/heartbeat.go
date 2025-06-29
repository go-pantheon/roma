package system

import (
	"time"

	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/mercury/gen/task/system"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"google.golang.org/protobuf/proto"
)

func NewHeartbeatTask() *system.HeartbeatTask {
	return system.NewHeartbeatTask(&climsg.CSHeartbeat{
		ClientTime: time.Now().Unix(),
	}, assertHeartBeatExec)
}

func assertHeartBeatExec(ctx core.Worker, cs, sc proto.Message) (err error) {
	return nil
}
