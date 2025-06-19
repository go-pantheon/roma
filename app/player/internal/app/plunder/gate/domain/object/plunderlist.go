package object

import (
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = "plunder"
)

func init() {
	userregister.Register(ModuleKey, NewPlunderList)
}

var _ life.Module = (*PlunderList)(nil)

type PlunderList struct {
	Plunders map[int64]*Plunder
}

func NewPlunderList() life.Module {
	o := &PlunderList{
		Plunders: make(map[int64]*Plunder),
	}

	return o
}

func (o *PlunderList) Marshal() ([]byte, error) {
	p := dbv1.PlunderListProtoPool.Get()
	defer dbv1.PlunderListProtoPool.Put(p)

	p.Plunders = make(map[int64]*dbv1.PlunderProto, len(o.Plunders))

	for k, op := range o.Plunders {
		p.Plunders[k] = NewPlunderProto()
		op.encodeServer(p.Plunders[k])
	}

	return proto.Marshal(p)
}

func (o *PlunderList) Unmarshal(data []byte) error {
	p := dbv1.PlunderListProtoPool.Get()
	defer dbv1.PlunderListProtoPool.Put(p)

	if err := proto.Unmarshal(data, p); err != nil {
		return errors.Wrap(err, "failed to unmarshal plunder list")
	}

	for k, p := range p.Plunders {
		o.Plunders[k] = NewPlunder()
		o.Plunders[k].decodeServer(p)
	}

	return nil
}
