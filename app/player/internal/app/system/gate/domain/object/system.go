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

var _ life.Module = (*System)(nil)

type System struct {
	firstHeartBeatCompleted atomic.Bool
	// TODO: add more fields
}

func NewSystem() *System {
	o := &System{}
	o.Register()

	return o
}

func (o *System) Register() {
	userregister.Register(ModuleKey, o)
}

func (o *System) Marshal() ([]byte, error) {
	p := dbv1.SystemProtoPool.Get()
	defer dbv1.SystemProtoPool.Put(p)

	return proto.Marshal(p)
}

func (o *System) Unmarshal(bytes []byte) error {
	p := dbv1.SystemProtoPool.Get()
	defer dbv1.SystemProtoPool.Put(p)

	if err := proto.Unmarshal(bytes, p); err != nil {
		return errors.Wrap(err, "failed to unmarshal system")
	}

	return nil
}

func (o *System) FirstHeartBeatCompleted() bool {
	return o.firstHeartBeatCompleted.Load()
}

func (o *System) SetFirstHeartBeatCompleted() {
	o.firstHeartBeatCompleted.Store(true)
}
