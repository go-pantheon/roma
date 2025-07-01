package object

import (
	"maps"
	"slices"

	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/zerrors"
)

type Hero struct {
	Id     int64
	Level  int64
	Skills map[int64]int64 // skill id -> level
	Equips []int64         // wearing equip ids
}

func NewHero(dataId int64) (*Hero, error) {
	d := gamedata.GetHeroBaseData(dataId)
	if d == nil {
		return nil, errors.Wrapf(zerrors.ErrGameDataNotFound, "id=%d", dataId)
	}

	r := &Hero{
		Id:     d.Id(),
		Level:  gamedata.DefaultHeroLevel,
		Skills: make(map[int64]int64, len(d.SkillDataList)),
	}

	return r, nil
}

func (o *Hero) decodeServer(p *dbv1.UserHeroProto) *Hero {
	if p == nil {
		return o
	}

	o.Level = p.Level
	o.Equips = slices.Clone(p.Equips)

	o.Skills = make(map[int64]int64, len(p.Skills))
	maps.Copy(o.Skills, p.Skills)

	return o
}

func (o *Hero) encodeServer() *dbv1.UserHeroProto {
	p := dbv1.UserHeroProtoPool.Get()

	p.Id = o.Id
	p.Level = o.Level
	p.Equips = slices.Clone(o.Equips)

	p.Skills = make(map[int64]int64, len(o.Skills))
	maps.Copy(p.Skills, o.Skills)

	return p
}

func (o *Hero) EncodeClient() *climsg.HeroProto {
	p := &climsg.HeroProto{
		Id:     o.Id,
		Level:  o.Level,
		Equips: slices.Clone(o.Equips),
		Skills: make(map[int64]int64, len(o.Skills)),
	}

	maps.Copy(p.Skills, o.Skills)

	return p
}
