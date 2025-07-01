package object

import (
	"time"

	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-util/errors"
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
	return &Dev{}
}

func (o *Dev) IsLifeModule() {}

func (o *Dev) EncodeServer() proto.Message {
	p := dbv1.UserDevProtoPool.Get()

	p.TimeOffset = int64(o.devTimeOffset)

	return p
}

func (o *Dev) DecodeServer(p proto.Message) error {
	if p == nil {
		return errors.New("dev decode server nil")
	}

	op, ok := p.(*dbv1.UserDevProto)
	if !ok {
		return errors.Errorf("dev decode server invalid type: %T", p)
	}

	o.devTimeOffset = time.Duration(op.TimeOffset)

	return nil
}

func (o *Dev) SetTimeOffset(dur time.Duration) bool {
	o.devTimeOffset += dur
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
