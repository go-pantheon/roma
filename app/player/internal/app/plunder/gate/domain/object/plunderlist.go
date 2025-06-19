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

func (os *PlunderList) EncodeServer() proto.Message {
	p := dbv1.UserPlunderListProtoPool.Get()
	p.Plunders = make(map[int64]*dbv1.UserPlunderProto, len(os.Plunders))

	for k, op := range os.Plunders {
		pp := dbv1.UserPlunderProtoPool.Get()

		op.encodeServer(pp)
		p.Plunders[k] = pp
	}

	return p
}

func (os *PlunderList) DecodeServer(p proto.Message) error {
	if p == nil {
		return errors.New("plunder list decode server nil")
	}

	sp, ok := p.(*dbv1.UserPlunderListProto)
	if !ok {
		return errors.Errorf("plunder list decode server invalid type: %T", p)
	}

	os.Plunders = make(map[int64]*Plunder, len(sp.Plunders))

	for id, op := range sp.Plunders {
		o := NewPlunder()

		o.decodeServer(op)
		os.Plunders[id] = o
	}

	return nil
}
