package domain

import (
	"github.com/go-kratos/kratos/log"
	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/hero/gate/domain/object"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	"github.com/vulcan-frame/vulcan-game/gamedata"
)

type HeroDomain struct {
	log *log.Helper
}

func NewHeroDomain(logger log.Logger) *HeroDomain {
	return &HeroDomain{
		log: log.NewHelper(log.With(logger, "module", "player/hero/gate/domain")),
	}
}

func (do *HeroDomain) UnlockHero(ctx core.Context, d *gamedata.HeroBaseData) (hero *object.Hero, err error) {
	heroList := ctx.User().HeroList
	if heroList.Heroes[d.ID] != nil {
		err = errors.Errorf("hero already unlocked. id=%d", d.ID)
		return
	}

	hero, err = object.NewHero(d.ID)
	if err != nil {
		return
	}

	hero.Level = d.Level
	heroList.Heroes[d.ID] = hero
	return
}
