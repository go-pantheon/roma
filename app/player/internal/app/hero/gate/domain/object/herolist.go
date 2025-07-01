package object

import (
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = life.ModuleKey("hero")
)

func init() {
	userregister.Register(ModuleKey, NewHeroList)
}

var _ life.Module = (*HeroList)(nil)

type HeroList struct {
	Heroes map[int64]*Hero
}

func NewHeroList() life.Module {
	o := &HeroList{
		Heroes: make(map[int64]*Hero, 8),
	}

	return o
}

func (o *HeroList) IsLifeModule() {}

func (o *HeroList) EncodeServer() proto.Message {
	p := dbv1.UserHeroListProtoPool.Get()
	p.Heroes = make(map[int64]*dbv1.UserHeroProto, len(o.Heroes))

	for _, oh := range o.Heroes {
		p.Heroes[oh.Id] = oh.encodeServer()
	}

	return p
}

func (o *HeroList) DecodeServer(p proto.Message) error {
	if p == nil {
		return errors.New("hero list decode server nil")
	}

	op, ok := p.(*dbv1.UserHeroListProto)
	if !ok {
		return errors.Errorf("hero list decode server invalid type: %T", p)
	}

	o.Heroes = make(map[int64]*Hero, len(op.Heroes))

	for _, ph := range op.Heroes {
		h, err := NewHero(ph.Id)
		if err != nil {
			return err
		}

		o.Heroes[ph.Id] = h.decodeServer(ph)
	}

	return nil
}

func (o *HeroList) EncodeClient() *climsg.UserHeroListProto {
	p := &climsg.UserHeroListProto{
		Heroes: make(map[int64]*climsg.HeroProto, len(o.Heroes)),
	}

	for _, oh := range o.Heroes {
		p.Heroes[oh.Id] = oh.EncodeClient()
	}

	return p
}
