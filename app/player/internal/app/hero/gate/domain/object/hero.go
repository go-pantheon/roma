package object

import (
	"context"
	"maps"
	"slices"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/gamedata"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
	"github.com/vulcan-frame/vulcan-game/pkg/errs"
)

const DefaultHeroLevel = 1

type HeroList struct {
	Heroes map[int64]*Hero
}

type Hero struct {
	Id     int64
	Level  int64
	Skills map[int64]int64 // skill id -> level
	Equips []int64         // wearing equip ids
}

func NewHeroList() *HeroList {
	return &HeroList{
		Heroes: make(map[int64]*Hero, 8),
	}
}

func NewHero(dataId int64) (*Hero, error) {
	d := gamedata.GetHeroBaseData(dataId)
	if d == nil {
		return nil, errors.Wrapf(errs.ErrGameDataNotFound, "id=%d", dataId)
	}
	r := &Hero{
		Id:     d.Id(),
		Level:  DefaultHeroLevel,
		Skills: make(map[int64]int64, len(d.SkillDataList)),
	}
	return r, nil
}

func (o *HeroList) DecodeServer(ctx context.Context, p *dbv1.UserHeroListProto) error {
	if p == nil {
		return nil
	}

	o.Heroes = make(map[int64]*Hero, len(p.Heroes))

	for _, ph := range p.Heroes {
		h, err := NewHero(ph.Id)
		if err != nil {
			return err
		}
		o.Heroes[ph.Id], err = h.DecodeServer(ctx, ph)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Hero) DecodeServer(ctx context.Context, p *dbv1.UserHeroProto) (*Hero, error) {
	if p == nil {
		return o, nil
	}
	o.Level = int64(p.Level)
	o.Equips = slices.Clone(p.Equips)
	o.Skills = make(map[int64]int64, len(p.Skills))
	maps.Copy(o.Skills, p.Skills)

	return o, nil
}

func NewHeroListProto() *dbv1.UserHeroListProto {
	p := &dbv1.UserHeroListProto{}
	p.Heroes = make(map[int64]*dbv1.UserHeroProto, 8)
	return p
}

func (o *HeroList) EncodeServer(p *dbv1.UserHeroListProto) {
	p.Heroes = make(map[int64]*dbv1.UserHeroProto, len(o.Heroes))

	for _, oh := range o.Heroes {
		p.Heroes[oh.Id] = oh.EncodeServer()
	}
}

func (o *Hero) EncodeServer() *dbv1.UserHeroProto {
	p := &dbv1.UserHeroProto{
		Id:     o.Id,
		Level:  o.Level,
		Skills: make(map[int64]int64, len(o.Skills)),
		Equips: slices.Clone(o.Equips),
	}
	maps.Copy(p.Skills, o.Skills)
	return p
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

func (o *Hero) EncodeClient() *climsg.HeroProto {
	p := &climsg.HeroProto{
		Id:     o.Id,
		Level:  o.Level,
		Skills: make(map[int64]int64, len(o.Skills)),
		Equips: slices.Clone(o.Equips),
	}
	maps.Copy(p.Skills, o.Skills)

	return p
}
