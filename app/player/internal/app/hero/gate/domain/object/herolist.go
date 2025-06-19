package object

import (
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = "hero"
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

func (o *HeroList) Marshal() ([]byte, error) {
	p := dbv1.UserHeroListProtoPool.Get()
	defer dbv1.UserHeroListProtoPool.Put(p)

	p.Heroes = make(map[int64]*dbv1.UserHeroProto, len(o.Heroes))
	for _, oh := range o.Heroes {
		ph := dbv1.UserHeroProtoPool.Get()
		oh.encodeServer(ph)
		p.Heroes[oh.Id] = ph
	}

	return proto.Marshal(p)
}

func (o *HeroList) Unmarshal(data []byte) error {
	p := dbv1.UserHeroListProtoPool.Get()
	defer dbv1.UserHeroListProtoPool.Put(p)

	o.Heroes = make(map[int64]*Hero, len(p.Heroes))

	for _, ph := range p.Heroes {
		h, err := NewHero(ph.Id)
		if err != nil {
			return err
		}

		o.Heroes[ph.Id], err = h.decodeServer(ph)
		if err != nil {
			return err
		}
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
