package object

import (
	"time"

	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = life.ModuleKey("dev")
)

var _ life.Module = (*Dev)(nil)

func init() {
	userregister.Register(ModuleKey, NewDev)
}

type Dev struct {
	devTimeOffset time.Duration
}

func NewDev() life.Module {
	o := &Dev{}
	return o
}

func NewDevProto() *dbv1.DevProto {
	p := &dbv1.DevProto{}
	return p
}

func (o *Dev) EncodeServer(p *dbv1.DevProto) {
	p.TimeOffset = int64(o.devTimeOffset)
}

func (o *Dev) Marshal() ([]byte, error) {
	p := dbv1.DevProtoPool.Get()
	defer dbv1.DevProtoPool.Put(p)

	p.TimeOffset = int64(o.devTimeOffset)
	return proto.Marshal(p)
}

func (o *Dev) DecodeServer(p *dbv1.DevProto) {
	if p == nil {
		return
	}

	o.devTimeOffset = time.Duration(p.TimeOffset)
}

func (o *Dev) Unmarshal(bytes []byte) error {
	p := dbv1.DevProtoPool.Get()
	defer dbv1.DevProtoPool.Put(p)

	return proto.Unmarshal(bytes, p)
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
