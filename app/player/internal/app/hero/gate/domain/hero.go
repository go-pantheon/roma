package domain

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/hero/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	"github.com/pkg/errors"
)

type HeroDomain struct {
	log *log.Helper
}

func NewHeroDomain(logger log.Logger) *HeroDomain {
	return &HeroDomain{
		log: log.NewHelper(log.With(logger, "module", "player/hero/gate/domain")),
	}
}

func (do *HeroDomain) UnlockHero(ctx core.Context, d *gamedata.HeroBaseData) (*object.Hero, error) {
	heroList := ctx.User().HeroList()
	if heroList.Heroes[d.ID] != nil {
		return nil, errors.Errorf("hero already unlocked. id=%d", d.ID)
	}

	hero, err := object.NewHero(d.ID)
	if err != nil {
		return nil, err
	}

	hero.Level = d.Level
	heroList.Heroes[d.ID] = hero

	return hero, nil
}
