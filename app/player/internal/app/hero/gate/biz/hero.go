package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/hero/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/hero/gate/domain/object"
	storagedo "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/pkg/errors"
)

func NewHeroUseCase(mgr *core.Manager, logger log.Logger,
	heroDo *domain.HeroDomain,
	storageDo *storagedo.StorageDomain,
) *HeroUseCase {
	uc := &HeroUseCase{
		log:       log.NewHelper(log.With(logger, "module", "player/hero/gate/biz")),
		mgr:       mgr,
		heroDo:    heroDo,
		storageDo: storageDo,
	}

	mgr.RegisterOnCreatedEvent(uc.onCreated)
	mgr.RegisterEvent(core.EventStorageItemUpdated, uc.OnStorageUpdated)
	return uc
}

type HeroUseCase struct {
	log       *log.Helper
	mgr       *core.Manager
	heroDo    *domain.HeroDomain
	storageDo *storagedo.StorageDomain
}

func (uc *HeroUseCase) onCreated(ctx core.Context) error {
	for _, d := range gamedata.GetHeroConstantData().HeroDataList {
		if _, err := uc.heroDo.UnlockHero(ctx, d); err != nil {
			uc.log.WithContext(ctx).Errorf("unlock hero on created failed. %+v", err)
			continue
		}
		ctx.Changed(object.ModuleKey)
	}

	return nil
}

func (uc *HeroUseCase) OnStorageUpdated(ctx core.Context, arg *life.EventArg) error {
	itemIds := life.GetArg[[]int64](arg, core.ArgKeyStorageItemIDs)

	var needUnlock bool
	for _, itemId := range itemIds {
		if itemData := gamedata.GetResourceItemData(itemId); itemData != nil {
			if itemData.Type == gamedata.ItemTypeHeroPiece {
				needUnlock = true
				break
			}
		}
	}
	if !needUnlock {
		return nil
	}

	newHeroes := make([]*object.Hero, 0, len(gamedata.GetHeroBaseDataList()))

	heroes := ctx.User().HeroList().Heroes
	for _, heroData := range gamedata.GetHeroBaseDataList() {
		if _, ok := heroes[heroData.ID]; ok {
			continue
		}
		if uc.storageDo.CanCost(ctx, heroData.CostCosts) == nil {
			if newHero, err := uc.unlockHero(ctx, heroData, heroData.CostCosts); err != nil {
				uc.log.WithContext(ctx).Errorf("unlock hero failed. %+v", err)
			} else {
				newHeroes = append(newHeroes, newHero)
			}
		}
	}

	msg := &climsg.SCPushHeroUnlock{
		Heroes: make([]*climsg.HeroProto, 0, len(newHeroes)),
	}
	for _, hero := range newHeroes {
		msg.Heroes = append(msg.Heroes, hero.EncodeClient())
	}
	_ = ctx.Reply(int32(climod.ModuleID_Hero), int32(cliseq.HeroSeq_PushHeroUnlock), ctx.UID(), msg)
	return nil
}

func (uc *HeroUseCase) unlockHero(ctx core.Context, d *gamedata.HeroBaseData, cost *gamedata.Costs) (r *object.Hero, err error) {
	if ctx.User().HeroList().Heroes[d.ID] != nil {
		err = errors.Errorf("hero already unlocked. heroId=%d", d.ID)
		return
	}

	if err = uc.storageDo.Cost(ctx, cost); err != nil {
		return
	}

	r, err = uc.heroDo.UnlockHero(ctx, d)
	if err != nil {
		return
	}

	ctx.Changed(object.ModuleKey)
	return
}

func (uc *HeroUseCase) HeroLevelUpgrade(ctx core.Context, cs *climsg.CSHeroLevelUpgrade) (sc *climsg.SCHeroLevelUpgrade, err error) {
	sc = &climsg.SCHeroLevelUpgrade{}
	var hero *object.Hero
	if hero = ctx.User().HeroList().Heroes[cs.HeroId]; hero == nil {
		sc.Code = climsg.SCHeroLevelUpgrade_ErrHeroNotExist
		return
	}
	heroData := gamedata.GetHeroBaseData(hero.Id)
	if heroData == nil {
		sc.Code = climsg.SCHeroLevelUpgrade_ErrHeroNotExist
		return
	}
	levelData := heroData.LevelData.SubDataMap[int64(hero.Level)+1]
	if levelData == nil {
		sc.Code = climsg.SCHeroLevelUpgrade_ErrMaxLevel
		return
	}
	if len(levelData.Cost) > 0 {
		if err = uc.storageDo.Cost(ctx, levelData.CostCosts); err != nil {
			sc.Code = climsg.SCHeroLevelUpgrade_ErrCostNotEnough
			return
		}
	}
	hero.Level = levelData.Level

	ctx.Changed(object.ModuleKey)

	sc.Code = climsg.SCHeroLevelUpgrade_Succeeded
	sc.Hero = hero.EncodeClient()
	return
}
