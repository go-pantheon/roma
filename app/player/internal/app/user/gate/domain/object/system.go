package userobj

import (
	"sync/atomic"

	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
)

type System struct {
	firstHeartBeatCompleted atomic.Bool
}

func NewSystem() *System {
	o := &System{}
	return o
}

func NewSystemProto() *dbv1.SystemProto {
	p := &dbv1.SystemProto{}
	return p
}

func (o *System) EncodeServer(p *dbv1.SystemProto) {
}

func (o *System) DecodeServer(p *dbv1.SystemProto) {
	if p == nil {
		return
	}
}

func (o *System) FirstHeartBeatCompleted() bool {
	return o.firstHeartBeatCompleted.Load()
}

func (o *System) SetFirstHeartBeatCompleted() {
	o.firstHeartBeatCompleted.Store(true)
}
