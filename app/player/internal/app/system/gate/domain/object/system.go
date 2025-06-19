package systemobj

import (
	"sync/atomic"

	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = "system"
)

func init() {
	userregister.Register(ModuleKey, NewSystem)
}

var _ life.Module = (*System)(nil)

type System struct {
	firstHeartBeatCompleted atomic.Bool
	// TODO: add more fields
}

func NewSystem() life.Module {
	return &System{}
}

func (o *System) EncodeServer() proto.Message {
	p := dbv1.UserSystemProtoPool.Get()

	return p
}

func (o *System) DecodeServer(p proto.Message) error {
	if p == nil {
		return errors.New("system decode server nil")
	}

	_, ok := p.(*dbv1.UserSystemProto)
	if !ok {
		return errors.Errorf("system decode server invalid type: %T", p)
	}

	return nil
}

func (o *System) FirstHeartBeatCompleted() bool {
	return o.firstHeartBeatCompleted.Load()
}

func (o *System) SetFirstHeartBeatCompleted() {
	o.firstHeartBeatCompleted.Store(true)
}
