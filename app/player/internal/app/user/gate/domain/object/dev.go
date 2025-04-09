package userobj

import (
	"time"

	"github.com/go-pantheon/fabrica-kit/profile"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

type Dev struct {
	devTimeOffset time.Duration
}

func NewDev() *Dev {
	return &Dev{}
}

func NewDevProto() *dbv1.DevProto {
	p := &dbv1.DevProto{}
	return p
}

func (o *Dev) EncodeServer(p *dbv1.DevProto) {
	p.TimeOffset = int64(o.devTimeOffset)
}

func (o *Dev) DecodeServer(p *dbv1.DevProto) {
	if p == nil {
		return
	}

	o.devTimeOffset = time.Duration(p.TimeOffset)
}

func (o *Dev) SetTimeOffset(dur time.Duration) bool {
	o.devTimeOffset = o.devTimeOffset + dur
	return true
}

func (o *Dev) ResetTimeOffset() {
	o.devTimeOffset = time.Duration(0)
}

func (o *Dev) TimeOffset() time.Duration {
	if profile.IsDev() {
		return o.devTimeOffset
	}
	return 0
}

func (o *User) Now() time.Time {
	if profile.IsDev() {
		return time.Now().Add(o.Dev.TimeOffset())
	}
	return time.Now()
}
